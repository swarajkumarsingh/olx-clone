package vault

import (
	"olx-clone/functions/general"
	"testing"
)

func TestEncryption(t *testing.T) {
	type testStruct struct {
		word string
	}
	var testCases = []testStruct{}
	// Add 10 Random string to the test cases
	for i := 0; i < 10; i++ {
		testCase, _ := general.GenerateRandomString(20)
		testCases = append(testCases, testStruct{testCase})
	}
	for index, testCase := range testCases {
		encrypted, err := AESEncrypt(testCase.word)
		if err != nil {
			// Fail the test
			t.Error("Failed to encrypt the text", err)
		}
		output, err := AESDecrypt(encrypted)
		if err != nil {
			// Fail the test
			t.Error("Failed to decrypt the text", err)
		}
		if testCase.word != output {
			t.Errorf("Case %d: Output %v not equal to expected %v", index+1, output, testCase.word)
		}
	}
}
