
priv/android-apk.keystore:
	mkdir -p priv
	keytool -genkey -v -keystore priv/android-apk.keystore -alias FyneCameraDemo -keyalg RSA -keysize 2048 -validity 36500

signer-sign-FyneCameraDemo-arm64-apk: priv/android-apk.keystore
	zipalign -p 4 ./examples/FyneCameraDemo/FyneCameraDemo.apk ./examples/FyneCameraDemo/FyneCameraDemo-aligned.apk
	rm -f ./examples/FyneCameraDemo/FyneCameraDemo.apk 
	mv ./examples/FyneCameraDemo/FyneCameraDemo-aligned.apk ./examples/FyneCameraDemo/FyneCameraDemo.apk
	zipalign -c 4 ./examples/FyneCameraDemo/FyneCameraDemo.apk
	apksigner sign --ks-key-alias FyneCameraDemo --ks priv/android-apk.keystore ./examples/FyneCameraDemo/FyneCameraDemo.apk || \
		apksigner sign --ks-key-alias FyneCameraDemo --ks priv/android-apk.keystore ./examples/FyneCameraDemo/FyneCameraDemo.apk || \
		apksigner sign --ks-key-alias FyneCameraDemo --ks priv/android-apk.keystore ./examples/FyneCameraDemo/FyneCameraDemo.apk || \
		apksigner sign --ks-key-alias FyneCameraDemo --ks priv/android-apk.keystore ./examples/FyneCameraDemo/FyneCameraDemo.apk

	apksigner verify ./examples/FyneCameraDemo/FyneCameraDemo.apk
