package main

import (
	"strings"
	"unicode"
)

// ToKatakana converts wapuro-hepburn romaji into the equivalent katakana.
func ToKatakana(s string) string {
	s = strings.ToUpper(s)
	s = preKatakana.Replace(s)
	s = romajiToKatakana.Replace(s)
	s = strings.Map(HiraganaToKatakana, s)
	s = postKatakana.Replace(s)
	s = postKanaSpecial.Replace(s)
	return removeNonJapaneseChars(s)
}

// preKatakana performs character transliterations before romajiToKatakana replacements
// have been performed. The order of calls is important to avoid aggressive replacement
// of elements within larger strings. The romaji keys are in uppercase to allow lower and
// upper case latin to convert to hiragana and katakana respectively, when then type-agnostic
// function ToKana is called.
var preKatakana = strings.NewReplacer(
	"XA", "ァ", // Replace wapruo x-prefixes with katakana small vowels.
	"XI", "ィ",
	"XU", "ゥ",
	"XE", "ェ",
	"XO", "ォ",

	"BB", "ッB", // Pre-convert double-consonants into っconsonant pairs for step 2 replacement.
	"TC", "ッC",
	"CC", "ッC",
	"DD", "ッD",
	"FF", "ッF",
	"GG", "ッG",
	"HH", "ッH",
	"JJ", "ッJ",
	"KK", "ッK",
	"LL", "ッL",
	"MM", "ッM",
	"PP", "ッP",
	"QQ", "ッQ",
	"RR", "ッR",
	"SS", "ッS",
	"TT", "ッT",
	"VV", "ッV",
	"WW", "ッW",
	"YY", "ッY",
	"ZZ", "ッZ",
)

// romajiToKatakana is a Replacer which maps romaji keys to katakana characters.
// In order to allow case sensitive replacements within ToKana, the romaji
// keys are all uppercase.
var romajiToKatakana = strings.NewReplacer(
	"KA", "カ",
	"KI", "キ",
	"KU", "ク",
	"KE", "ケ",
	"KO", "コ",
	"SA", "サ",
	"SHI", "シ",
	"SU", "ス",
	"SE", "セ",
	"SO", "ソ",
	"TA", "タ",
	"CHI", "チ",
	"TSU", "ツ",
	"TE", "テ",
	"TO", "ト",
	"NA", "ナ",
	"NI", "ニ",
	"NU", "ヌ",
	"NE", "ネ",
	"NO", "ノ",
	"HA", "ハ",
	"HI", "ヒ",
	"FU", "フ",
	"HE", "ヘ",
	"HO", "ホ",
	"MA", "マ",
	"MI", "ミ",
	"MU", "ム",
	"ME", "メ",
	"MO", "モ",
	"YA", "ヤ",
	"YU", "ユ",
	"YO", "ヨ",
	"RA", "ラ",
	"RI", "リ",
	"RU", "ル",
	"RE", "レ",
	"RO", "ロ",
	"WA", "ワ",
	"WO", "ヲ",
	"GA", "ガ",
	"GI", "ギ",
	"GU", "グ",
	"GE", "ゲ",
	"GO", "ゴ",
	"ZA", "ザ",
	"JI", "ジ",
	"ZU", "ズ",
	"ZE", "ゼ",
	"ZO", "ゾ",
	"DA", "ダ",
	"DI", "ヂ",
	"DU", "ヅ",
	"DE", "デ",
	"DO", "ド",
	"BA", "バ",
	"BI", "ビ",
	"BU", "ブ",
	"BE", "ベ",
	"BO", "ボ",
	"PA", "パ",
	"PI", "ピ",
	"PU", "プ",
	"PE", "ペ",
	"PO", "ポ",
	"DEU", "デュ",
	"KYA", "キャ",
	"KYU", "キュ",
	"KYO", "キョ",
	"SHA", "シャ",
	"SHU", "シュ",
	"SHO", "ショ",
	"CHA", "チャ",
	"CHU", "チュ",
	"CHO", "チョ",
	"NYA", "ニャ",
	"NYU", "ニュ",
	"NYO", "ニョ",
	"HYA", "ヒャ",
	"HYU", "ヒュ",
	"HYO", "ヒョ",
	"MYA", "ミャ",
	"MYU", "ミュ",
	"MYO", "ミョ",
	"RYA", "リャ",
	"RYU", "リュ",
	"RYO", "リョ",
	"GYA", "ギャ",
	"GYU", "ギュ",
	"GYO", "ギョ",
	"JA", "ジャ",
	"JU", "ジュ",
	"JO", "ジョ",
	"JYA", "ジャ",
	"JYU", "ジュ",
	"JYO", "ジョ",
	"DYA", "ヂャ",
	"DYU", "ヂュ",
	"DYO", "ヂョ",
	"BYA", "ビャ",
	"BYU", "ビュ",
	"BYO", "ビョ",
	"PYA", "ピャ",
	"PYU", "ピュ",
	"PYO", "ピョ",
	"YI", "イィ",
	"YE", "イェ",
	"WI", "ウィ",
	"WU", "ウゥ",
	"WE", "ウェ",
	"WYU", "ウュ",
	"VA", "ヴァ",
	"VI", "ヴィ",
	"VU", "ヴ",
	"VE", "ヴェ",
	"VO", "ヴォ",
	"VYA", "ヴャ",
	"VYU", "ヴュ",
	"VYE", "ヴィェ",
	"VYO", "ヴョ",
	"KYE", "キェ",
	"GYE", "ギェ",
	"KWA", "クァ",
	"KWI", "クィ",
	"KWE", "クェ",
	"KWU", "クゥ",
	"KWO", "クォ",
	"GWA", "グァ",
	"GWI", "グィ",
	"GWE", "グェ",
	"GWO", "グォ",
	"GWU", "グゥ",
	"SHE", "シェ",
	"JE", "ジェ",
	"SI", "スィ",
	"ZI", "ズィ",
	"CHE", "チェ",
	"TSA", "ツァ",
	"TSE", "ツェ",
	"TSI", "ツィ",
	"TSO", "ツォ",
	"TSYU", "ツュ",
	"TI", "ティ",
	"TU", "トゥ",
	"TYU", "テュ",
	"NYE", "ニェ",
	"HYE", "ヒェ",
	"BYE", "ビェ",
	"PYE", "ピェ",
	"FA", "ファ",
	"FI", "フィ",
	"FE", "フェ",
	"FO", "フォ",
	"FYA", "フャ",
	"FYU", "フュ",
	"FYE", "フィェ",
	"FYO", "フョ",
	"HU", "ホゥ",
	"MYE", "ミェ",
	"RYE", "リェ",
	"LA", "ラ",
	"LI", "リ",
	"LU", "ル",
	"LE", "レ",
	"LO", "ロ",
	"QA", "クァ",
	"QI", "クィ",
	"QE", "クェ",
	"QO", "クォ",
	"QU", "クヮ",
	"A", "ア",
	"I", "イ",
	"U", "ウ",
	"E", "エ",
	"O", "オ",
	"N", "ン",
)

// HiraganaToKatakana replaces a single hiragana character with the
// unicode equivalent katakana character.
func HiraganaToKatakana(r rune) rune {
	if (r >= 'ぁ' && r <= 'ゖ') || (r >= 'ゝ' && r <= 'ゞ') {
		return r + 0x60
	}
	return r
}

// postKatakana performs final character transliterations after romajiToKatakana
// replacements have occurred.
var postKatakana = strings.NewReplacer(
	"X", "ッ", // any dangling wapruo x-prefixes become katakana small tu (0x30C3.
)

// postKanaSpecial performs final character transliterations after all others have
// been performed.
var postKanaSpecial = strings.NewReplacer(
	"–", "ー", // convert en-dash (0x2013) to katakana-hiragana prolonged sound mark (0x30FC).
	"-", "ー", // convert hyphen-minus (0x2D) to (0x30FC).
	"'", "", // strip out single quotes used to designated moriac n's.
)

func removeNonJapaneseChars(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Han) {
			result = append(result, r)
		}
	}
	return string(result)
}
