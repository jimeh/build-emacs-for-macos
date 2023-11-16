package dmgbuild

import (
	"fmt"
	"sort"
	"strings"
)

type locale string

//nolint:golint
var (
	LocaleAfZA     locale = "af_ZA"
	LocaleAr       locale = "ar"
	LocaleBeBY     locale = "be_BY"
	LocaleBgBG     locale = "bg_BG"
	LocaleBn       locale = "bn"
	LocaleBo       locale = "bo"
	LocaleBr       locale = "br"
	LocaleCaES     locale = "ca_ES"
	LocaleCsCZ     locale = "cs_CZ"
	LocaleCy       locale = "cy"
	LocaleDaDK     locale = "da_DK"
	LocaleDeAT     locale = "de_AT"
	LocaleDeCH     locale = "de_CH"
	LocaleDeDE     locale = "de_DE"
	LocaleDzBT     locale = "dz_BT"
	LocaleElCY     locale = "el_CY"
	LocaleElGR     locale = "el_GR"
	LocaleEnAU     locale = "en_AU"
	LocaleEnCA     locale = "en_CA"
	LocaleEnGB     locale = "en_GB"
	LocaleEnIE     locale = "en_IE"
	LocaleEnSG     locale = "en_SG"
	LocaleEnUS     locale = "en_US"
	LocaleEo       locale = "eo"
	LocaleEs419    locale = "es_419"
	LocaleEsES     locale = "es_ES"
	LocaleEtEE     locale = "et_EE"
	LocaleFaIR     locale = "fa_IR"
	LocaleFiFI     locale = "fi_FI"
	LocaleFoFO     locale = "fo_FO"
	LocaleFr001    locale = "fr_001"
	LocaleFrBE     locale = "fr_BE"
	LocaleFrCA     locale = "fr_CA"
	LocaleFrCH     locale = "fr_CH"
	LocaleFrFR     locale = "fr_FR"
	LocaleGaLatgIE locale = "ga-Latg_IE"
	LocaleGaIE     locale = "ga_IE"
	LocaleGd       locale = "gd"
	LocaleGrc      locale = "grc"
	LocaleGuIN     locale = "gu_IN"
	LocaleGv       locale = "gv"
	LocaleHeIL     locale = "he_IL"
	LocaleHiIN     locale = "hi_IN"
	LocaleHrHR     locale = "hr_HR"
	LocaleHuHU     locale = "hu_HU"
	LocaleHyAM     locale = "hy_AM"
	LocaleIsIS     locale = "is_IS"
	LocaleItCH     locale = "it_CH"
	LocaleItIT     locale = "it_IT"
	LocaleIuCA     locale = "iu_CA"
	LocaleJaJP     locale = "ja_JP"
	LocaleKaGE     locale = "ka_GE"
	LocaleKl       locale = "kl"
	LocaleKoKR     locale = "ko_KR"
	LocaleLtLT     locale = "lt_LT"
	LocaleLvLV     locale = "lv_LV"
	LocaleMkMK     locale = "mk_MK"
	LocaleMrIN     locale = "mr_IN"
	LocaleMtMT     locale = "mt_MT"
	LocaleNbNO     locale = "nb_NO"
	LocaleNeNP     locale = "ne_NP"
	LocaleNlBE     locale = "nl_BE"
	LocaleNlNL     locale = "nl_NL"
	LocaleNnNO     locale = "nn_NO"
	LocalePa       locale = "pa"
	LocalePlPL     locale = "pl_PL"
	LocalePtBR     locale = "pt_BR"
	LocalePtPT     locale = "pt_PT"
	LocaleRoRO     locale = "ro_RO"
	LocaleRuRU     locale = "ru_RU"
	LocaleSe       locale = "se"
	LocaleSkSK     locale = "sk_SK"
	LocaleSlSI     locale = "sl_SI"
	LocaleSrRS     locale = "sr_RS"
	LocaleSvSE     locale = "sv_SE"
	LocaleThTH     locale = "th_TH"
	LocaleToTO     locale = "to_TO"
	LocaleTrTR     locale = "tr_TR"
	LocaleUkUA     locale = "uk_UA"
	LocaleUrIN     locale = "ur_IN"
	LocaleUrPK     locale = "ur_PK"
	LocaleUzUZ     locale = "uz_UZ"
	LocaleViVN     locale = "vi_VN"
	LocaleZhCN     locale = "zh_CN"
	LocaleZhTW     locale = "zh_TW"
)

type Buttons struct {
	LanguageName string
	Agree        string
	Disagree     string
	Print        string
	Save         string
	Message      string
}

type License struct {
	DefaultLanguage locale
	Licenses        map[locale]string
	Buttons         map[locale]Buttons
}

func NewLicense() License {
	return License{}
}

//nolint:goconst
func (s *License) Render() []string {
	var l []string

	if s.DefaultLanguage != "" {
		l = append(l,
			"\"default-language\": "+pyStr(string(s.DefaultLanguage)),
		)
	}

	if len(s.Licenses) > 0 {
		var items []string
		for k, v := range s.Licenses {
			items = append(items, fmt.Sprintf(
				"%s: %s", pyStr(string(k)), pyMStr(v),
			))
		}
		sort.SliceStable(items, func(i, j int) bool {
			return items[i] < items[j]
		})
		l = append(l,
			"\"licenses\": {\n        "+
				strings.Join(items, ",\n        ")+
				"\n    }",
		)
	}

	if len(s.Buttons) > 0 {
		var items []string
		for k, v := range s.Buttons {
			items = append(items, fmt.Sprintf(
				"%s: (\n"+
					"            %s,\n"+
					"            %s,\n"+
					"            %s,\n"+
					"            %s,\n"+
					"            %s,\n"+
					"            %s\n"+
					"        )",
				pyStr(string(k)),
				pyStr(v.LanguageName),
				pyStr(v.Agree),
				pyStr(v.Disagree),
				pyStr(v.Print),
				pyStr(v.Save),
				pyStr(v.Message),
			))
		}
		sort.SliceStable(items, func(i, j int) bool {
			return items[i] < items[j]
		})
		l = append(l,
			"\"buttons\": {\n        "+
				strings.Join(items, ",\n        ")+
				"\n    }",
		)
	}

	if len(l) == 0 {
		return []string{}
	}

	return []string{
		"license = {\n    " + strings.Join(l, ",\n    ") + "\n}\n",
	}
}
