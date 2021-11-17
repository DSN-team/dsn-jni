package main

// #include <stdlib.h>
// #include <stddef.h>
// #include <stdint.h>
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/ClarkGuan/jni"
	"github.com/DSN-team/core"
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

//export Java_com_dsnteam_dsn_CoreManager_loadProfiles
func Java_com_dsnteam_dsn_CoreManager_loadProfiles(env uintptr, _ uintptr) int {
	return core.LoadProfiles()
}

//export Java_com_dsnteam_dsn_CoreManager_getProfilesIds
func Java_com_dsnteam_dsn_CoreManager_getProfilesIds(env uintptr, _ uintptr) (ids uintptr) {
	ids = jni.Env(env).NewIntArray(len(core.Profiles))
	for i := 0; i < len(core.Profiles); i++ {
		jni.Env(env).SetIntArrayElement(ids, i, core.Profiles[i].Id)
	}
	return
}

//export Java_com_dsnteam_dsn_CoreManager_getProfilesNames
func Java_com_dsnteam_dsn_CoreManager_getProfilesNames(env uintptr, _ uintptr) (usernames uintptr) {
	dataType := jni.Env(env).FindClass("Ljava/lang/String;")
	usernames = jni.Env(env).NewObjectArray(len(core.Profiles), dataType, 0)
	for i := 0; i < len(core.Profiles); i++ {
		jni.Env(env).SetObjectArrayElement(usernames, i, jni.Env(env).NewString(core.Profiles[i].Username))
	}
	return
}

//export Java_com_dsnteam_dsn_CoreManager_getProfilePublicKey
func Java_com_dsnteam_dsn_CoreManager_getProfilePublicKey(env uintptr, _ uintptr) uintptr {
	friend := Friend{Username: currentProfile.ThisUser.Username, Address: currentProfile.ThisUser.Address, PublicKey: currentProfile.GetProfilePublicKey()}
	b, err := json.Marshal(friend)
	core.ErrHandler(err)
	return jni.Env(env).NewString(string(b))
}

//export Java_com_dsnteam_dsn_CoreManager_getProfileName
func Java_com_dsnteam_dsn_CoreManager_getProfileName(env uintptr, _ uintptr) uintptr {
	return jni.Env(env).NewString(currentProfile.ThisUser.Username)
}

//export Java_com_dsnteam_dsn_CoreManager_getProfileAddress
func Java_com_dsnteam_dsn_CoreManager_getProfileAddress(env uintptr, _ uintptr) uintptr {
	return jni.Env(env).NewString(currentProfile.ThisUser.Address)
}

type Friend struct {
	Username  string
	Address   string
	PublicKey string
}

//export Java_com_dsnteam_dsn_CoreManager_addFriend
func Java_com_dsnteam_dsn_CoreManager_addFriend(env uintptr, _ uintptr, dataIn uintptr) {
	data := string(jni.Env(env).GetStringUTF(dataIn))
	var friend Friend
	err := json.Unmarshal([]byte(data), &friend)
	core.ErrHandler(err)
	currentProfile.AddFriend(friend.Username, friend.Address, friend.PublicKey)
}

//export Java_com_dsnteam_dsn_CoreManager_loadFriends
func Java_com_dsnteam_dsn_CoreManager_loadFriends(env uintptr, _ uintptr) int {
	return currentProfile.LoadFriends()
}

//export Java_com_dsnteam_dsn_CoreManager_getFriendsIds
func Java_com_dsnteam_dsn_CoreManager_getFriendsIds(env uintptr, _ uintptr) (ids uintptr) {
	ids = jni.Env(env).NewIntArray(len(currentProfile.Friends))
	for i := 0; i < len(currentProfile.Friends); i++ {
		jni.Env(env).SetIntArrayElement(ids, i, currentProfile.Friends[i].Id)
	}
	return
}

//export Java_com_dsnteam_dsn_CoreManager_getFriendsNames
func Java_com_dsnteam_dsn_CoreManager_getFriendsNames(env uintptr, _ uintptr) (usernames uintptr) {
	//friends = getFriends()
	dataType := jni.Env(env).FindClass("Ljava/lang/String;")
	usernames = jni.Env(env).NewObjectArray(len(currentProfile.Friends), dataType, 0)
	for i := 0; i < len(currentProfile.Friends); i++ {
		jni.Env(env).SetObjectArrayElement(usernames, i, jni.Env(env).NewString(currentProfile.Friends[i].Username))
	}
	return
}

//export Java_com_dsnteam_dsn_CoreManager_getFriendsAddresses
func Java_com_dsnteam_dsn_CoreManager_getFriendsAddresses(env uintptr, _ uintptr) (address uintptr) {
	//friends = getFriends()
	dataType := jni.Env(env).FindClass("Ljava/lang/String;")
	address = jni.Env(env).NewObjectArray(len(currentProfile.Friends), dataType, 0)
	for i := 0; i < len(currentProfile.Friends); i++ {
		jni.Env(env).SetObjectArrayElement(address, i, jni.Env(env).NewString(currentProfile.Friends[i].Address))
	}
	return
}

//export Java_com_dsnteam_dsn_CoreManager_getFriendsPublicKeys
func Java_com_dsnteam_dsn_CoreManager_getFriendsPublicKeys(env uintptr, _ uintptr) (publicKey uintptr) {
	//friends = getFriends()
	dataType := jni.Env(env).FindClass("Ljava/lang/String;")
	publicKey = jni.Env(env).NewObjectArray(len(currentProfile.Friends), dataType, 0)
	for i := 0; i < len(currentProfile.Friends); i++ {
		jni.Env(env).SetObjectArrayElement(publicKey, i, jni.Env(env).NewString(core.EncPublicKey(core.MarshalPublicKey(currentProfile.Friends[i].PublicKey))))
	}
	return
}

//export Java_com_dsnteam_dsn_CoreManager_connectToFriends
func Java_com_dsnteam_dsn_CoreManager_connectToFriends(env uintptr, _ uintptr) {
	currentProfile.ConnectToFriends()
}

//export Java_com_dsnteam_dsn_CoreManager_connectToFriend
func Java_com_dsnteam_dsn_CoreManager_connectToFriend(env uintptr, _ uintptr, userId int) {
	currentProfile.ConnectToFriend(userId)
}

//export Java_com_dsnteam_dsn_CoreManager_getFriendsRequests
func Java_com_dsnteam_dsn_CoreManager_getFriendsRequests(env uintptr, _ uintptr) (usernames uintptr) {
	//friends = getFriends()
	dataType := jni.Env(env).FindClass("Ljava/lang/String;")
	usernames = jni.Env(env).NewObjectArray(len(currentProfile.FriendRequests), dataType, 0)
	for i := 0; i < len(currentProfile.FriendRequests); i++ {
		jni.Env(env).SetObjectArrayElement(usernames, i, jni.Env(env).NewString(currentProfile.FriendRequests[i].Username))
	}
	return
}

//Инициализировать структуры и подключение

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

	currentProfile.WriteRequest(
		userId,
		currentProfile.BuildDataRequest(core.RequestData,
			uint64(lenIn),
			currentProfile.DataStrInput,
			userId))
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
