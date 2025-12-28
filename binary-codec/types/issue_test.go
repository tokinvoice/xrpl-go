package types

import (
	"errors"
	"testing"

	"github.com/Peersyst/xrpl-go/binary-codec/types/testutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestIssue_FromJson(t *testing.T) {
	tt := []struct {
		name        string
		input       any
		expected    []byte
		expectedErr error
	}{
		{
			name: "pass - valid xrp issue object",
			input: map[string]any{
				"currency": "XRP",
			},
			expected: []byte{
				0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			expectedErr: nil,
		},
		{
			name: "pass - valid issue iou object",
			input: map[string]any{
				"currency": "USD",
				"issuer":   "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn",
			},
			expected: []byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 85, 83, 68, 0, 0, 0, 0, 0,
				174, 18, 58, 133, 86, 243, 207, 145, 21, 71,
				17, 55, 106, 251, 15, 137, 79, 131, 43, 61,
			},
		},
		{
			name: "pass - valid xrp issue object",
			input: map[string]any{
				"currency": "0123456789ABCDEF0123456789ABCDEF01234567",
				"issuer":   "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn",
			},
			expected: []byte{
				1, 35, 69, 103, 137, 171, 205, 239, 1,
				35, 69, 103, 137, 171, 205, 239, 1, 35,
				69, 103, 174, 18, 58, 133, 86, 243, 207,
				145, 21, 71, 17, 55, 106, 251, 15, 137,
				79, 131, 43, 61,
			},
		},
		{
			name: "pass - valid mpt issuance id",
			input: map[string]any{
				"mpt_issuance_id": "BAADF00DBAADF00DBAADF00DBAADF00DBAADF00DBAADF00D",
			},
			expected: []byte{
				186,
				173,
				240,
				13,
				186,
				173,
				240,
				13,
				186,
				173,
				240,
				13,
				186,
				173,
				240,
				13,
				186,
				173,
				240,
				13,
				186,
				173,
				240,
				13,
			},
		},
		{
			name:        "fail - invalid Issue",
			input:       "r3e7qTG44Mg8pHXgxPtyRx286Re5Urtx2p2",
			expected:    nil,
			expectedErr: ErrInvalidIssueObject,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			issue := &Issue{}
			actual, err := issue.FromJSON(tc.input)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestIssue_ToJson(t *testing.T) {
	tt := []struct {
		name     string
		expected any
		opts     []int
		err      error
		setup    func(t *testing.T) (*Issue, *testutil.MockBinaryParser)
	}{
		{
			name: "pass - valid issue object",
			expected: map[string]any{
				"currency": "USD",
				"issuer":   "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn",
			},
			opts: []int{20},
			err:  nil,
			setup: func(t *testing.T) (*Issue, *testutil.MockBinaryParser) {
				ctrl := gomock.NewController(t)
				mock := testutil.NewMockBinaryParser(ctrl)
				mock.EXPECT().ReadBytes(20).Return([]byte{
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 85, 83, 68, 0, 0, 0, 0, 0,
				}, nil)
				mock.EXPECT().ReadBytes(20).Return([]byte{
					174, 18, 58, 133, 86, 243, 207, 145, 21, 71,
					17, 55, 106, 251, 15, 137, 79, 131, 43, 61,
				}, nil)
				return &Issue{}, mock
			},
		},
		{
			name: "pass - valid xrp issue object",
			expected: map[string]any{
				"currency": "XRP",
			},
			opts: []int{20},
			err:  nil,
			setup: func(t *testing.T) (*Issue, *testutil.MockBinaryParser) {
				ctrl := gomock.NewController(t)
				mock := testutil.NewMockBinaryParser(ctrl)
				mock.EXPECT().ReadBytes(20).Return(XRPBytes, nil)
				return &Issue{}, mock
			},
		},
		{
			name: "pass - mpt issuance id",
			expected: map[string]any{
				"mpt_issuance_id": "BAADF00DBAADF00DBAADF00DBAADF00DBAADF00DBAADF00D",
			},
			opts: []int{MPTIssuanceIDBytesLength},
			err:  nil,
			setup: func(t *testing.T) (*Issue, *testutil.MockBinaryParser) {
				ctrl := gomock.NewController(t)
				mock := testutil.NewMockBinaryParser(ctrl)
				mock.EXPECT().ReadBytes(MPTIssuanceIDBytesLength).Return([]byte{
					186,
					173,
					240,
					13,
					186,
					173,
					240,
					13,
					186,
					173,
					240,
					13,
					186,
					173,
					240,
					13,
					186,
					173,
					240,
					13,
					186,
					173,
					240,
					13,
				}, nil)
				return &Issue{}, mock
			},
		},
		{
			name:     "fail - invalid Issue",
			expected: nil,
			opts:     []int{20},
			err:      errors.New("errReadBytes"),
			setup: func(t *testing.T) (*Issue, *testutil.MockBinaryParser) {
				ctrl := gomock.NewController(t)
				mock := testutil.NewMockBinaryParser(ctrl)
				mock.EXPECT().ReadBytes(20).Return([]byte{}, errors.New("errReadBytes"))
				return &Issue{}, mock
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			issue, parser := tc.setup(t)
			actual, err := issue.ToJSON(parser, tc.opts...)

			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

// TestIssue_isIssueObject tests the isIssueObject helper function.
func TestIssue_isIssueObject(t *testing.T) {
	issue := &Issue{}

	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{
			name:     "nil",
			input:    nil,
			expected: false,
		},
		{
			name:     "string",
			input:    "not a map",
			expected: false,
		},
		{
			name:     "empty map",
			input:    map[string]any{},
			expected: false,
		},
		{
			name:     "currency only (XRP)",
			input:    map[string]any{"currency": "XRP"},
			expected: true,
		},
		{
			name:     "mpt_issuance_id only",
			input:    map[string]any{"mpt_issuance_id": "00000001B5F762798A53D543A014CAF8B297CFF8F2F937E8"},
			expected: true,
		},
		{
			name:     "currency and issuer (IOU)",
			input:    map[string]any{"currency": "USD", "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"},
			expected: true,
		},
		{
			name:     "wrong single key",
			input:    map[string]any{"foo": "bar"},
			expected: false,
		},
		{
			name:     "currency with wrong second key",
			input:    map[string]any{"currency": "USD", "foo": "bar"},
			expected: false,
		},
		{
			name:     "three keys",
			input:    map[string]any{"currency": "USD", "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh", "extra": "key"},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := issue.isIssueObject(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

// TestIssue_FromJSON_InvalidInputs tests error handling for invalid inputs.
func TestIssue_FromJSON_InvalidInputs(t *testing.T) {
	issue := &Issue{}

	tests := []struct {
		name    string
		input   any
		wantErr error
	}{
		{
			name:    "nil input",
			input:   nil,
			wantErr: ErrInvalidIssueObject,
		},
		{
			name:    "string input",
			input:   "not a map",
			wantErr: ErrInvalidIssueObject,
		},
		{
			name:    "integer input",
			input:   123,
			wantErr: ErrInvalidIssueObject,
		},
		{
			name:    "slice input",
			input:   []string{"a", "b"},
			wantErr: ErrInvalidIssueObject,
		},
		{
			name:    "empty map",
			input:   map[string]any{},
			wantErr: ErrInvalidIssueObject,
		},
		{
			name:    "map with wrong keys",
			input:   map[string]any{"foo": "bar"},
			wantErr: ErrInvalidIssueObject,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := issue.FromJSON(tc.input)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

// TestIssue_FromJSON_InvalidMPTIssuanceID tests error handling for invalid MPT issuance IDs.
func TestIssue_FromJSON_InvalidMPTIssuanceID(t *testing.T) {
	issue := &Issue{}

	tests := []struct {
		name  string
		input map[string]any
	}{
		{
			name:  "non-hex mpt_issuance_id",
			input: map[string]any{"mpt_issuance_id": "not-hex-string"},
		},
		{
			name:  "odd-length hex mpt_issuance_id",
			input: map[string]any{"mpt_issuance_id": "ABC"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := issue.FromJSON(tc.input)
			require.Error(t, err)
		})
	}
}
