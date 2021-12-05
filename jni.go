package main

// #include <stdlib.h>
// #include <stddef.h>
// #include <stdint.h>
import "C"
import (
	"fmt"
	"github.com/ClarkGuan/jni"
	"github.com/DSN-team/core"
	utils2 "github.com/DSN-team/core/utils"
	"log"
	"reflect"
	"runtime"
	"unsafe"
)

var workingVM jni.VM
var callBackBufferPtr unsafe.Pointer
var callBackBufferCap int
var currentProfile core.Profile

func main() {
	fmt.Println("test")
}

//export Java_com_dsnteam_dsn_CoreManager_initDB
func Java_com_dsnteam_dsn_CoreManager_initDB(env uintptr, _ uintptr) {
	core.StartDB()
}

//export Java_com_dsnteam_dsn_CoreManager_register
func Java_com_dsnteam_dsn_CoreManager_register(env uintptr, _ uintptr, usernameIn uintptr, passwordIn uintptr, addressIn uintptr) (result bool) {
	username, password, address := string(jni.Env(env).GetStringUTF(usernameIn)), string(jni.Env(env).GetStringUTF(passwordIn)), string(jni.Env(env).GetStringUTF(addressIn))
	return currentProfile.Register(username, password, address)
}

//export Java_com_dsnteam_dsn_CoreManager_login
func Java_com_dsnteam_dsn_CoreManager_login(env uintptr, _ uintptr, pos int, passwordIn uintptr) (result bool) {
	password := string(jni.Env(env).GetStringUTF(passwordIn))
	return currentProfile.Login(password, pos)
}

//export Java_com_dsnteam_dsn_CoreManager_runServer
func Java_com_dsnteam_dsn_CoreManager_runServer(env uintptr, _ uintptr, addressIn uintptr) {
	address := string(jni.Env(env).GetStringUTF(addressIn))
	println("env run:", env)
	if env != 0 {
		workingVM, _ = jni.Env(env).GetJavaVM()
	}
	core.UpdateUI = updateCall
	currentProfile.RunServer(address)
}

//export Java_com_dsnteam_dsn_CoreManager_setCallBackBuffer
func Java_com_dsnteam_dsn_CoreManager_setCallBackBuffer(env uintptr, _ uintptr, jniBuffer uintptr) {
	callBackBufferPtr = jni.Env(env).GetDirectBufferAddress(jniBuffer)
	callBackBufferCap = jni.Env(env).GetDirectBufferCapacity(jniBuffer)
}

//export Java_com_dsnteam_dsn_CoreManager_writeBytes
func Java_com_dsnteam_dsn_CoreManager_writeBytes(env uintptr, _ uintptr, inBuffer uintptr, lenIn int, userId int) {
	log.Println("env write:", env)
	defer runtime.KeepAlive(currentProfile.DataStrInput)
	point, size := jni.Env(env).GetDirectBufferAddress(inBuffer), jni.Env(env).GetDirectBufferCapacity(inBuffer)

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&currentProfile.DataStrInput))
	sh.Data = uintptr(point)
	sh.Len = lenIn
	sh.Cap = size

	userPos, _ := currentProfile.FriendsIDXs.Load(uint(userId))

	dataMessageEncrypted := currentProfile.BuildDataMessage(currentProfile.DataStrInput, uint(userId))
	request := core.Request{RequestType: utils2.RequestData, PublicKey: core.MarshalPublicKey(&currentProfile.PrivateKey.PublicKey), Data: dataMessageEncrypted}

	currentProfile.WriteRequest(currentProfile.Friends[userPos.(int)], request)
}

//Realisation for platform
func updateCall(count int, userId int) {
	//Call Application to read structure and update internal data interpretations, update UI.
	var env jni.Env
	env, _ = workingVM.AttachCurrentThread()
	println("WorkingEnv:", env)
	classInput := env.FindClass("com/dsnteam/dsn/CoreManager")
	println("class_input:", classInput)
	methodId := env.GetStaticMethodID(classInput, "getUpdateCallBack", "(II)V")
	println("MethodID:", methodId)
	var bData []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&bData))
	sh.Data = uintptr(callBackBufferPtr)
	sh.Cap = callBackBufferCap
	sh.Len = callBackBufferCap
	println("buffer pointer:", callBackBufferPtr)
	copy(bData, currentProfile.DataStrOutput)
	println("buffer write done")
	env.CallStaticVoidMethodA(classInput, methodId, jni.Jvalue(count), jni.Jvalue(userId))
	workingVM.DetachCurrentThread()
	runtime.KeepAlive(bData)
}
