package models

import (
	"testing"
)

func TestCheckNewPost(t *testing.T) {
	tests := []struct {
		name     string
		postData PostData
		wantMsgs int
	}{
		{
			name: "Valid Post",
			postData: PostData{
				Title:      "Valid Title",
				Content:    "Valid Content",
				Categories: []string{"Technology"},
			},
			wantMsgs: 0,
		},
		{
			name: "Missing Categories",
			postData: PostData{
				Title:      "Valid Title",
				Content:    "Valid Content",
				Categories: []string{},
			},
			wantMsgs: 1,
		},
		{
			name: "Too Many Categories",
			postData: PostData{
				Title:      "Valid Title",
				Content:    "Valid Content",
				Categories: []string{"Cat1", "Cat2", "Cat3", "Cat4"},
			},
			wantMsgs: 1,
		},
		{
			name: "Empty Title",
			postData: PostData{
				Title:      "",
				Content:    "Valid Content",
				Categories: []string{"Technology"},
			},
			wantMsgs: 1,
		},
		{
			name: "Title Too Long",
			postData: PostData{
				Title:      "This title is definitely more than seventy characters long just to make sure the validation works correctly",
				Content:    "Valid Content",
				Categories: []string{"Technology"},
			},
			wantMsgs: 1,
		},
		{
			name: "Content Too Long",
			postData: PostData{
				Title:      "Valid Title",
				Content:    string(make([]byte, 5001)),
				Categories: []string{"Technology"},
			},
			wantMsgs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := CheckNewPost(tt.postData)
			if len(response.Messages) != tt.wantMsgs {
				t.Errorf("CheckNewPost() got %v messages, want %v. Msgs: %v", len(response.Messages), tt.wantMsgs, response.Messages)
			}
		})
	}
}
