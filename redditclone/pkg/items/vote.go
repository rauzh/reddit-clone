package items

type Vote struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

func (item *Item) NewVote(vote int, userID string) {
	vt := &Vote{
		Vote: vote,
		User: userID,
	}

	item.Votes = append(item.Votes, vt)
	item.Score += vote

	upvotes := (len(item.Votes) + item.Score) / 2
	item.UpvotePercentage = int(float32(upvotes) / float32(len(item.Votes)) * 100)
}

func (item *Item) DeleteVote(userID string) (bool, error) {
	i := -1
	vt := 0
	for idx, vote := range item.Votes {
		if vote.User == userID {
			i = idx
			vt = vote.Vote
			break
		}
	}
	if i == -1 {
		return false, nil
	}

	if i < len(item.Votes)-1 {
		copy(item.Votes[i:], item.Votes[i+1:])
	}
	item.Votes[len(item.Votes)-1] = nil
	item.Votes = item.Votes[:len(item.Votes)-1]

	item.Score -= vt
	upvotes := (len(item.Votes) + item.Score) / 2
	if len(item.Votes) != 0 {
		item.UpvotePercentage = int((float32(upvotes) / float32(len(item.Votes))) * 100)
	} else {
		item.UpvotePercentage = 0
	}
	return true, nil
}
