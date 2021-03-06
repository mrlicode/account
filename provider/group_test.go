package provider

import (
	"context"
	"testing"

	"github.com/grpcbrick/account/standard"
)

func TestService_CreateGroup(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.CreateGroupRequest
		wantState standard.State
		wantErr   bool
	}{
		{"正常创建", &standard.CreateGroupRequest{Name: "TEST", Category: "Category", State: "State", Description: "Description"},
			standard.State_SUCCESS, false},
		{"正常创建", &standard.CreateGroupRequest{Name: "TEST2", Category: "Category", State: "State", Description: "Description"},
			standard.State_SUCCESS, false},
		{"正常创建", &standard.CreateGroupRequest{Name: "TEST3", Category: "Category", State: "State", Description: "Description"},
			standard.State_SUCCESS, false},
		{"重复的 Name", &standard.CreateGroupRequest{Name: "TEST", Category: "Category", State: "State", Description: "Description"},
			standard.State_GROUP_ALREADY_EXISTS, false},
		{"空的 Name", &standard.CreateGroupRequest{Name: "", Category: "Category", State: "Nickname", Description: "Username"},
			standard.State_PARAMS_INVALID, false},
		{"空的 Category", &standard.CreateGroupRequest{Name: "TEST", Category: "", State: "Nickname", Description: "Username"},
			standard.State_PARAMS_INVALID, false},
		{"空的 State", &standard.CreateGroupRequest{Name: "TEST", Category: "Category", State: "", Description: "Username"},
			standard.State_PARAMS_INVALID, false},
		{"空的 Description", &standard.CreateGroupRequest{Name: "TEST", Category: "Category", State: "Nickname", Description: ""},
			standard.State_PARAMS_INVALID, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.CreateGroup(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.CreateGroup() = %v, want %v", gotResp, tt.wantState)
				return
			}

			if gotResp.State == standard.State_SUCCESS {
				if gotResp.Data.Name != tt.args.Name {
					t.Errorf("Service.CreateGroup() = %v, want %v", gotResp, tt.wantState)
					return
				}
			}
		})
	}
}

func TestService_QueryGroups(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name            string
		args            *standard.QueryGroupsRequest
		wantState       standard.State
		wantDataSize    int64
		wantCurrentPage int64
		wantTotalPage   int64
		wantErr         bool
	}{
		{"正常查询", &standard.QueryGroupsRequest{Page: 1, Limit: 90},
			standard.State_SUCCESS, 3, 1, 1, false},
		{"只查一条", &standard.QueryGroupsRequest{Page: 1, Limit: 1},
			standard.State_SUCCESS, 1, 1, 3, false},
		{"第二页", &standard.QueryGroupsRequest{Page: 1, Limit: 2},
			standard.State_SUCCESS, 2, 1, 2, false},
		{"空的 ID", &standard.QueryGroupsRequest{Page: 1, Limit: 0},
			standard.State_PARAMS_INVALID, 0, 1, 0, false},
		{"不存在的 ID", &standard.QueryGroupsRequest{},
			standard.State_PARAMS_INVALID, 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.QueryGroups(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.QueryGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.QueryGroups() = %v, want %v", gotResp, tt.wantState)
				return
			}

			if tt.wantState == standard.State_SUCCESS {
				if gotResp.TotalPage != tt.wantTotalPage {
					t.Errorf("Service.QueryGroups() = %v, want %d", gotResp, tt.wantTotalPage)
					return
				}
				if int64(len(gotResp.Data)) != tt.wantDataSize {
					t.Errorf("Service.QueryGroups() = %v, want %d", gotResp.Data, tt.wantDataSize)
					return
				}
				if gotResp.CurrentPage != tt.wantCurrentPage {
					t.Errorf("Service.QueryUsers() = %v, want %d", gotResp, tt.wantCurrentPage)
					return
				}
			}
		})
	}
}

func TestService_QueryGroupByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.QueryGroupByIDRequest
		wantState standard.State
		wantName  string
		wantErr   bool
	}{
		{"正常查询", &standard.QueryGroupByIDRequest{ID: 1},
			standard.State_SUCCESS, "TEST", false},
		{"空的 ID", &standard.QueryGroupByIDRequest{},
			standard.State_PARAMS_INVALID, "ignore", false},
		{"不存在的 ID", &standard.QueryGroupByIDRequest{ID: 9999999},
			standard.State_GROUP_NOT_EXIST, "ignore", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.QueryGroupByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.QueryGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.QueryGroupByID() = %v, want %v", gotResp, tt.wantState)
				return
			}

			if tt.wantName != "ignore" {
				if gotResp.Data.Name != tt.wantName {
					t.Errorf("Service.QueryGroupByID() = %v, want %v", gotResp.Data.Name, tt.wantName)
				}
			}
		})
	}
}

func TestService_DeleteGroupByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.DeleteGroupByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.DeleteGroupByIDRequest{},
			standard.State_PARAMS_INVALID, false},
		{"正常删除", &standard.DeleteGroupByIDRequest{ID: 1},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.DeleteGroupByIDRequest{ID: 9999999},
			standard.State_GROUP_NOT_EXIST, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.DeleteGroupByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DeleteGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.DeleteGroupByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateGroupNameByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateGroupNameByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateGroupNameByIDRequest{Name: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateGroupNameByIDRequest{ID: 2, Name: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateGroupNameByIDRequest{ID: 99999, Name: "Update1"},
			standard.State_GROUP_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateGroupNameByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateGroupNameByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateGroupNameByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateGroupNameByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateGroupCategoryByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateGroupCategoryByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateGroupCategoryByIDRequest{Category: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateGroupCategoryByIDRequest{ID: 2, Category: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateGroupCategoryByIDRequest{ID: 99999, Category: "Update1"},
			standard.State_GROUP_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateGroupCategoryByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateGroupCategoryByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateGroupCategoryByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateGroupCategoryByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateGroupStateByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateGroupStateByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateGroupStateByIDRequest{State: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateGroupStateByIDRequest{ID: 2, State: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateGroupStateByIDRequest{ID: 99999, State: "Update1"},
			standard.State_GROUP_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateGroupStateByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateGroupStateByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateGroupStateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateGroupStateByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateGroupDescriptionByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateGroupDescriptionByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateGroupDescriptionByIDRequest{Description: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateGroupDescriptionByIDRequest{ID: 2, Description: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateGroupDescriptionByIDRequest{ID: 99999, Description: "Update1"},
			standard.State_GROUP_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateGroupDescriptionByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateGroupDescriptionByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateGroupDescriptionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateGroupDescriptionByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_AddUserToGroupByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.AddUserToGroupByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的用户 ID", &standard.AddUserToGroupByIDRequest{GroupID: 1},
			standard.State_PARAMS_INVALID, false},
		{"不存在的用户 ID", &standard.AddUserToGroupByIDRequest{ID: 9999, GroupID: 1},
			standard.State_USER_NOT_EXIST, false},
		{"空的组 ID", &standard.AddUserToGroupByIDRequest{ID: 1},
			standard.State_PARAMS_INVALID, false},
		{"不存在的组 ID", &standard.AddUserToGroupByIDRequest{ID: 1, GroupID: 999999},
			standard.State_GROUP_NOT_EXIST, false},
		{"正常添加", &standard.AddUserToGroupByIDRequest{ID: 1, GroupID: 1},
			standard.State_SUCCESS, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.AddUserToGroupByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddUserToGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.AddUserToGroupByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_RemoveUserFromGroupByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.RemoveUserFromGroupByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的用户 ID", &standard.RemoveUserFromGroupByIDRequest{ID: 1},
			standard.State_PARAMS_INVALID, false},
		{"不存在的用户 ID", &standard.RemoveUserFromGroupByIDRequest{UserID: 999999, ID: 1},
			standard.State_USER_NOT_EXIST, false},
		{"空的组 ID", &standard.RemoveUserFromGroupByIDRequest{UserID: 1},
			standard.State_PARAMS_INVALID, false},
		{"不存在的组 ID", &standard.RemoveUserFromGroupByIDRequest{UserID: 1, ID: 999999},
			standard.State_GROUP_NOT_EXIST, false},
		{"正常移除", &standard.RemoveUserFromGroupByIDRequest{UserID: 1, ID: 1},
			standard.State_SUCCESS, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.RemoveUserFromGroupByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RemoveUserFromGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.RemoveUserFromGroupByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}
