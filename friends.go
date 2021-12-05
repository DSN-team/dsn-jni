package main

import (
	"encoding/json"
	"github.com/ClarkGuan/jni"
	"github.com/DSN-team/core"
)

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
		jni.Env(env).SetIntArrayElement(ids, i, int(currentProfile.Friends[i].ID))
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
		jni.Env(env).SetObjectArrayElement(publicKey, i, jni.Env(env).NewString(core.EncodeKey(core.MarshalPublicKey(currentProfile.Friends[i].PublicKey))))
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
	currentProfile.FriendRequests = currentProfile.GetFriendRequests()
	dataType := jni.Env(env).FindClass("Ljava/lang/String;")
	usernames = jni.Env(env).NewObjectArray(len(currentProfile.FriendRequests), dataType, 0)
	for i := 0; i < len(currentProfile.FriendRequests); i++ {
		user := currentProfile.GetUser(currentProfile.FriendRequests[i].UserID)
		jni.Env(env).SetObjectArrayElement(usernames, i, jni.Env(env).NewString(user.Username))
	}
	return
}

//export Java_com_dsnteam_dsn_CoreManager_acceptFriendRequest
func Java_com_dsnteam_dsn_CoreManager_acceptFriendRequest(env uintptr, _ uintptr, pos int) {
	currentProfile.FriendRequests[pos].Accept()
}

//export Java_com_dsnteam_dsn_CoreManager_rejectFriendRequest
func Java_com_dsnteam_dsn_CoreManager_rejecttFriendRequest(env uintptr, _ uintptr, pos int) {
	currentProfile.FriendRequests[pos].Reject()
}
