package service

import "sort"

func fromAPIRequest(req PackRequest) (uint, error) {
	if req.Items == 0 {
		return 0, ErrEmptyItems
	}

	return req.Items, nil
}

func toAPIResponse(boxes []uint) PackResponse {
	var resp PackResponse

	orderMap := make(map[uint]uint)
	for i := range boxes {
		orderMap[boxes[i]]++
	}

	for k, v := range orderMap {
		resp.Packs = append(resp.Packs, Pack{
			Box:      k,
			Quantity: v,
		})
	}

	sort.Slice(resp.Packs, func(i, j int) bool {
		return resp.Packs[i].Box > resp.Packs[j].Box
	})

	return resp
}
