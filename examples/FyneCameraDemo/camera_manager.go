package main

import (
	"container/list"
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/camera"
	"github.com/xaionaro-go/ndk/util"
)

type CameraState struct {
	bAvailable bool
	devices    *list.List
}

func (s *CameraState) AddDevice(device *camera.Device) {
	if s.devices == nil {
		s.devices = list.New()
	}
	s.devices.PushBack(device)
}

func (s *CameraState) DelDevice(device *camera.Device) {
	util.Assert(s.devices != nil)
	elm := s.devices.Front()
	for ; elm != nil; elm = elm.Next() {
		if device == elm.Value.(*camera.Device) {
			s.devices.Remove(elm)
			return
		}
	}
}

func (s *CameraState) String() string {
	ret := ""
	for elm := s.devices.Front(); elm != nil; elm = elm.Next() {
		ret += fmt.Sprint(elm.Value, ",")
	}
	if ret != "" {
		return fmt.Sprint("{ ", s.bAvailable, ", [ ", ret, " ] }")
	} else {
		return fmt.Sprint("{ ", s.bAvailable, " }")
	}
}

type CameraManager struct {
	*camera.Manager
	cameras map[string]*CameraState
}

var gMgr *CameraManager

func CameraManagerInstance() *CameraManager {
	if gMgr != nil {
		return gMgr
	}

	mgr := &CameraManager{}
	mgr.Manager = camera.ManagerCreate()
	if mgr.Manager != nil {
		mgr.Manager.RegisterAvailabilityCallback(mgr)
		mgr.cameras = map[string]*CameraState{}

		gMgr = mgr
	}
	return gMgr
}

func (mgr *CameraManager) Delete() {
	mgr.Manager.UnregisterAvailabilityCallback(mgr)
	mgr.Manager.Delete()
	mgr.Manager = nil
}

func (mgr *CameraManager) cameraState(id string) *CameraState {
	if ret, ok := mgr.cameras[id]; ok {
		return ret
	}
	ret := new(CameraState)
	mgr.cameras[id] = ret
	return ret
}

func (mgr *CameraManager) OnCameraAvailable(id string) {
	log.Println("... onCameraAvailable", id)
	mgr.cameraState(id).bAvailable = true
}

func (mgr *CameraManager) OnCameraUnavailable(id string) {
	log.Println("... onCameraUnavailable", id)
	mgr.cameraState(id).bAvailable = false
}

func (mgr *CameraManager) OnDisconnected(device *camera.Device) {
	log.Println("... OnDisconnected", device)
	mgr.cameraState(device.GetId()).DelDevice(device)
}

func (mgr *CameraManager) OnError(device *camera.Device, error int) {
	log.Println("... OnError", device, error)
	mgr.cameraState(device.GetId()).DelDevice(device)
}

func (mgr *CameraManager) dump() string {
	cstr := ""
	for id, s := range mgr.cameras {
		cstr += id + " = " + s.String() + ";"
	}
	if cstr == "" {
		return ""
	}
	return cstr
}

func (mgr *CameraManager) GetSensorOrientation(id string) (int, int) {
	metadata, err := mgr.GetCameraCharacteristics(id)
	util.Assert(err)
	defer metadata.Free()

	lens, err := metadata.GetConstEntry(camera.LENS_FACING)
	util.Assert(err)
	orientation, err := metadata.GetConstEntry(camera.SENSOR_ORIENTATION)
	util.Assert(err)

	return int(lens.Data().([]byte)[0]), int(orientation.Data().([]int32)[0])
}

func (mgr *CameraManager) GetSupportPixels(id string) [][2]int {
	metadata, err := mgr.GetCameraCharacteristics(id)
	util.Assert(err)
	defer metadata.Free()

	cfgs, err := metadata.GetConstEntry(camera.SCALER_AVAILABLE_STREAM_CONFIGURATIONS)
	util.Assert(err)

	// SCALER_AVAILABLE_STREAM_CONFIGURATIONS:
	//  input = entry.data.i32[i * 4 + 3];
	//  format = entry.data.i32[i * 4 + 0];
	//  width = entry.data.i32[i * 4 + 1];
	//  height = entry.data.i32[i * 4 + 2];

	ps := [][2]int{}
	ds := cfgs.Data().([]int32)
	// Remove duplicates
	ms := map[string]bool{}
	for i := 0; i < cfgs.Count()/4; i++ {
		w, h := int(ds[i*4+1]), int(ds[i*4+2])
		s := fmt.Sprint(w, "x", h)
		if !ms[s] {
			ps = append(ps, [2]int{w, h})
			ms[s] = true
		}
	}
	return ps
}

func (mgr *CameraManager) OpenCamera(
	cameraID string,
) (*Camera, error) {
	dev, err := mgr.Manager.OpenCamera(cameraID, mgr)
	if err != nil {
		return nil, fmt.Errorf("unable to open the camera: %w", err)
	}

	mgr.cameraState(cameraID).AddDevice(dev)

	return &Camera{
		Device: dev,
	}, nil
}
