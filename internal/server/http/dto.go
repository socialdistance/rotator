package internalhttp

type ErrorDto struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type BannerToSlotDto struct {
	BannerID int64 `json:"banner_id"`
	SlotID   int64 `json:"slot_id"`
}

type CountTransitionDto struct {
	BannerID      int64 `json:"banner_id"`
	SlotID        int64 `json:"slot_id"`
	SocialGroupID int64 `json:"social_group_id"`
}

type ChooseBannerDto struct {
	SlotID        int64 `json:"slot_id"`
	SocialGroupID int64 `json:"social_group_id"`
}
