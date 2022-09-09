package hebcal

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
// Derived from original C version, Copyright (C) 1994-2004 Danny Sadinoff
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

import (
	"errors"
	"strconv"
	"time"
)

// DafYomi represents a page of Talmud, such as Pesachim 103
type DafYomi struct {
	Name  string // Tractate name (e..g Berachot)
	Blatt int    // Page number
}

var shas0 = []DafYomi{
	{Name: "Berachot", Blatt: 64},
	{Name: "Shabbat", Blatt: 157},
	{Name: "Eruvin", Blatt: 105},
	{Name: "Pesachim", Blatt: 121},
	{Name: "Shekalim", Blatt: 22},
	{Name: "Yoma", Blatt: 88},
	{Name: "Sukkah", Blatt: 56},
	{Name: "Beitzah", Blatt: 40},
	{Name: "Rosh Hashana", Blatt: 35},
	{Name: "Taanit", Blatt: 31},
	{Name: "Megillah", Blatt: 32},
	{Name: "Moed Katan", Blatt: 29},
	{Name: "Chagigah", Blatt: 27},
	{Name: "Yevamot", Blatt: 122},
	{Name: "Ketubot", Blatt: 112},
	{Name: "Nedarim", Blatt: 91},
	{Name: "Nazir", Blatt: 66},
	{Name: "Sotah", Blatt: 49},
	{Name: "Gitin", Blatt: 90},
	{Name: "Kiddushin", Blatt: 82},
	{Name: "Baba Kamma", Blatt: 119},
	{Name: "Baba Metzia", Blatt: 119},
	{Name: "Baba Batra", Blatt: 176},
	{Name: "Sanhedrin", Blatt: 113},
	{Name: "Makkot", Blatt: 24},
	{Name: "Shevuot", Blatt: 49},
	{Name: "Avodah Zarah", Blatt: 76},
	{Name: "Horayot", Blatt: 14},
	{Name: "Zevachim", Blatt: 120},
	{Name: "Menachot", Blatt: 110},
	{Name: "Chullin", Blatt: 142},
	{Name: "Bechorot", Blatt: 61},
	{Name: "Arachin", Blatt: 34},
	{Name: "Temurah", Blatt: 34},
	{Name: "Keritot", Blatt: 28},
	{Name: "Meilah", Blatt: 22},
	{Name: "Kinnim", Blatt: 4},
	{Name: "Tamid", Blatt: 9},
	{Name: "Midot", Blatt: 5},
	{Name: "Niddah", Blatt: 73},
}

var osday, nsday int

// GetDafYomi calculates the Daf Yomi for given date.
//
// Returns an error if the date is before Daf Yomi cycle began
// (Hebrew year 5684, 11 September 1923).
func GetDafYomi(date HDate) (DafYomi, error) {
	if osday == 0 {
		osday, _ = GregorianToRD(1923, time.September, 11)
		nsday, _ = GregorianToRD(1975, time.June, 24)
	}

	cday := date.Abs()
	if cday < osday {
		return DafYomi{}, errors.New("before Daf Yomi cycle began")
	}

	var cno, dno int
	if cday >= nsday { // "new" cycle
		cno = 8 + ((cday - nsday) / 2711)
		dno = (cday - nsday) % 2711
	} else { // old cycle
		cno = 1 + ((cday - osday) / 2702)
		dno = (cday - osday) % 2702
	}

	// Find the daf taking note that the cycle changed slightly after cycle 7.

	var total = 0
	var blatt = 0
	var count = -1

	var shas = shas0
	// Fix Shekalim for old cycles
	if cno <= 7 {
		shas = make([]DafYomi, len(shas0))
		copy(shas, shas0)
		shas[4] = DafYomi{Name: "Shekalim", Blatt: 13}
	}

	// Find the daf
	var dafcnt = 40
	for j := 0; j < dafcnt; j++ {
		count++
		total = total + shas[j].Blatt - 1
		if dno < total {
			blatt = (shas[j].Blatt + 1) - (total - dno)
			// fiddle with the weird ones near the end
			switch count {
			case 36:
				blatt = blatt + 21
			case 37:
				blatt = blatt + 24
			case 38:
				blatt = blatt + 32
			default:
				break
			}
			// Bailout
			j = 1 + dafcnt
		}
	}
	return DafYomi{Name: shas[count].Name, Blatt: blatt}, nil
}

// Returns a string representation of the Daf Yomi
func (daf DafYomi) String() string {
	return daf.Name + " " + strconv.Itoa(daf.Blatt)
}

type dafYomiEvent struct {
	Date HDate
	Daf  DafYomi
}

func (ev dafYomiEvent) GetDate() HDate {
	return ev.Date
}

func (ev dafYomiEvent) Render() string {
	return ev.Daf.String()
}

func (ev dafYomiEvent) GetFlags() HolidayFlags {
	return DAF_YOMI
}

func (ev dafYomiEvent) GetEmoji() string {
	return ""
}
