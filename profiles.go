package main

import (
	"encoding/json"
	"github.com/ClarkGuan/jni"
	"github.com/DSN-team/core"
)

//export Java_com_dsnteam_dsn_CoreManager_loadProfiles
func Java_com_dsnteam_dsn_CoreManager_loadProfiles(env uintptr, _ uintptr) int {
	return core.LoadProfiles()
}

//export Java_com_dsnteam_dsn_CoreManager_getProfilesIds
func Java_com_dsnteam_dsn_CoreManager_getProfilesIds(env uintptr, _ uintptr) (ids uintptr) {
	ids = jni.Env(env).NewIntArray(len(core.Profiles))
	for i := 0; i < len(core.Profiles); i++ {
		jni.Env(env).SetIntArrayElement(ids, i, int(core.Profiles[i].ID))
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
	friend := Friend{Username: currentProfile.Username, Address: currentProfile.Address, PublicKey: currentProfile.GetProfilePublicKey()}
	b, err := json.Marshal(friend)
	core.ErrHandler(err)
	return jni.Env(env).NewString(string(b))
}

//export Java_com_dsnteam_dsn_CoreManager_getProfileName
func Java_com_dsnteam_dsn_CoreManager_getProfileName(env uintptr, _ uintptr) uintptr {
	return jni.Env(env).NewString(currentProfile.Username)
}

//export Java_com_dsnteam_dsn_CoreManager_getProfileAddress
func Java_com_dsnteam_dsn_CoreManager_getProfileAddress(env uintptr, _ uintptr) uintptr {
	return jni.Env(env).NewString(currentProfile.Address)
}
