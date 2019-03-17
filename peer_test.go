package bcache

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaveworks/mesh"
)

func TestPeerOnGossip(t *testing.T) {
	testCases := []struct {
		name    string
		initial map[string]entry
		newMsg  map[string]entry
		delta   map[string]entry
	}{
		{
			name:    "from empty",
			initial: map[string]entry{},
			newMsg: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
			delta: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
		},
		{
			name: "new key",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
			delta: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
		},
		{
			name: "same key diff val",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 2,
				},
			},
			delta: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 2,
				},
			},
		},
		{
			name: "same key same val",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			delta: nil,
		},
		{
			name: "same key dif val same exp",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 1,
				},
			},
			delta: nil,
		},
	}

	var (
		peerID1 = mesh.PeerName(1)
		maxKeys = 100
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := newPeer(peerID1, maxKeys, &nopLogger{})
			require.NoError(t, err)

			// initial
			msg := newMessageFromEntries(peerID1, tc.initial)
			p.cc.mergeComplete(msg)

			// newMsg
			newMsg := newMessageFromEntries(peerID1, tc.newMsg)
			buf := newMsg.Encode()[0]
			delta, err := p.OnGossip(buf)
			require.NoError(t, err)

			if tc.delta == nil {
				require.Nil(t, delta)
			} else {
				require.Equal(t, tc.delta, delta.(*message).Entries)
			}
		})
	}
}

func TestPeerOnGossipBroadcast(t *testing.T) {
	testCases := []struct {
		name    string
		initial map[string]entry
		newMsg  map[string]entry
		delta   map[string]entry
	}{
		{
			name:    "from empty",
			initial: map[string]entry{},
			newMsg: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
			delta: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
		},
		{
			name: "new key",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
			delta: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
		},
		{
			name: "same key diff val",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 2,
				},
			},
			delta: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 2,
				},
			},
		},
		{
			name: "same key same val",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			delta: map[string]entry{}, // // OnGossipBroadcast returns received, which should never be nil
		},
		{
			name: "same key dif val same exp",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 1,
				},
			},
			delta: map[string]entry{}, // // OnGossipBroadcast returns received, which should never be nil
		},
	}

	var (
		peerID1 = mesh.PeerName(1)
		maxKeys = 100
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := newPeer(peerID1, maxKeys, &nopLogger{})
			require.NoError(t, err)

			// initial
			msg := newMessageFromEntries(peerID1, tc.initial)
			p.cc.mergeComplete(msg)

			// newMsg
			newMsg := newMessageFromEntries(peerID1, tc.newMsg)
			buf := newMsg.Encode()[0]
			delta, err := p.OnGossipBroadcast(mesh.UnknownPeerName, buf)
			require.NoError(t, err)

			require.Equal(t, tc.delta, delta.(*message).Entries)
		})
	}
}

func TestPeerOnGossipUnicast(t *testing.T) {
	testCases := []struct {
		name     string
		initial  map[string]entry
		newMsg   map[string]entry
		complete map[string]entry
	}{
		{
			name:    "from empty",
			initial: map[string]entry{},
			newMsg: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
			complete: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
		},
		{
			name: "new key",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
			complete: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
				"key2": entry{
					Key:     "key2",
					Val:     "val2",
					Expired: 1,
				},
			},
		},
		{
			name: "same key diff val",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 2,
				},
			},
			complete: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 2,
				},
			},
		},
		{
			name: "same key same val",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			complete: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
		},
		{
			name: "same key dif val same exp",
			initial: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
			newMsg: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val2",
					Expired: 1,
				},
			},
			complete: map[string]entry{
				"key1": entry{
					Key:     "key1",
					Val:     "val1",
					Expired: 1,
				},
			},
		},
	}

	var (
		peerID1 = mesh.PeerName(1)
		maxKeys = 100
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := newPeer(peerID1, maxKeys, &nopLogger{})
			require.NoError(t, err)

			// initial
			msg := newMessageFromEntries(peerID1, tc.initial)
			p.cc.mergeComplete(msg)

			// newMsg
			newMsg := newMessageFromEntries(peerID1, tc.newMsg)
			buf := newMsg.Encode()[0]
			err = p.OnGossipUnicast(mesh.UnknownPeerName, buf)
			require.NoError(t, err)

			require.Equal(t, tc.complete, p.cc.Messages().Entries)
		})
	}
}