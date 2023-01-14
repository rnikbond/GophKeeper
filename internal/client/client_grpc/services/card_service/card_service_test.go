package card_service

import (
	"GophKeeper/internal/model/card"
	"testing"
)

func Test_checkCardData(t *testing.T) {

	tests := []struct {
		name    string
		in      card.DataCardFull
		waitErr error
	}{
		{
			name: "Success",
			in: card.DataCardFull{
				MetaInfo: "MirPay",
				Number:   "4648289760410976",
				Period:   "10.2030",
				CVV:      "111",
				FullName: "Test Test",
			},
			waitErr: nil,
		},
		{
			name: "Check invalid number",
			in: card.DataCardFull{
				MetaInfo: "MirPay",
				Number:   "464289760410976",
				Period:   "10.2030",
				CVV:      "111",
				FullName: "Test Test",
			},
			waitErr: ErrInvalidNumber,
		},
		{
			name: "Check invalid period",
			in: card.DataCardFull{
				MetaInfo: "MirPay",
				Number:   "4648289760410976",
				Period:   "102030",
				CVV:      "111",
				FullName: "Test Test",
			},
			waitErr: ErrInvalidPeriod,
		},
		{
			name: "Check invalid CVV: chars",
			in: card.DataCardFull{
				MetaInfo: "MirPay",
				Number:   "4648289760410976",
				Period:   "10.2030",
				CVV:      "a1a",
				FullName: "Test Test",
			},
			waitErr: ErrInvalidCVV,
		},
		{
			name: "Check invalid CVV: short len",
			in: card.DataCardFull{
				MetaInfo: "MirPay",
				Number:   "4648289760410976",
				Period:   "10.2030",
				CVV:      "1",
				FullName: "Test Test",
			},
			waitErr: ErrInvalidCVV,
		},
		{
			name: "Check invalid CVV: short -12",
			in: card.DataCardFull{
				MetaInfo: "GPay",
				Number:   "4648289760410976",
				Period:   "10.2030",
				CVV:      "-12",
				FullName: "Test Test",
			},
			waitErr: ErrInvalidCVV,
		},
		{
			name: "Check invalid full name",
			in: card.DataCardFull{
				MetaInfo: "MirPay",
				Number:   "4648289760410976",
				Period:   "10.2030",
				CVV:      "111",
				FullName: "T",
			},
			waitErr: ErrInvalidFullName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := checkCardData(tt.in)

			if tt.waitErr != err {
				t.Errorf("checkCardData() error = %v, wantErr %v", err, tt.waitErr)
			}
		})
	}
}
