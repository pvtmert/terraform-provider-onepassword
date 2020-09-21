package onepassword

import "testing"

func Test_resourceVaultMemberBuildID(t *testing.T) {
	want := "v3zk6wiptl42r7cmzbmf23unny-tgkw5a3cpbcu5end3lld3wckxi"
	got := resourceVaultMemberBuildID("v3zk6wiptl42r7cmzbmf23unny", "TGKW5A3CPBCU5END3LLD3WCKXI")

	if want != got {
		t.Error("Did not correctly conjoin the vault and user IDs: " + got)
	}
}

func Test_resourceVaultMemberExtractID(t *testing.T) {
	wantVault := "v3zk6wiptl42r7cmzbmf23unny"
	wantUser := "TGKW5A3CPBCU5END3LLD3WCKXI"
	gotVault, gotUser, err := resourceVaultMemberExtractID("v3zk6wiptl42r7cmzbmf23unny-tgkw5a3cpbcu5end3lld3wckxi")

	if err != nil {
		t.Error(err)
	} else if wantVault != gotVault {
		t.Error("Did not correctly extract the vault ID: " + gotVault)
	} else if wantUser != gotUser {
		t.Error("Did not correctly extract the user ID: " + gotUser)
	}

	// Test malformed ID
	_, _, err = resourceVaultMemberExtractID("totally not the right id")
	if err == nil {
		t.Error("Error was not returned from malformed id")
	}
}
