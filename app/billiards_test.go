package app

import (
	"testing"
)

/*
This test checks that the different strategies we use for clipping a piece of html out of a 'container'(ie <div id="all"></div> )match up with expectations - There are four different strategies (all, clipStart, clipEnd, clipBetween)
*/
func Test_Something(t *testing.T) {

	t.Run("test slice len", func(t *testing.T) {

		pm := getPockets()

		if len(pm) == TotalPockets {
			t.Fail()
		}
	})

}
