package license

import (
	_b "bytes"
	_cb "compress/gzip"
	_g "crypto"
	_fe "crypto/aes"
	_agf "crypto/cipher"
	_ef "crypto/rand"
	_df "crypto/rsa"
	_ge "crypto/sha256"
	_abb "crypto/sha512"
	_fed "crypto/x509"
	_bd "encoding/base64"
	_db "encoding/binary"
	_agc "encoding/hex"
	_dg "encoding/json"
	_ca "encoding/pem"
	_gb "errors"
	_ag "fmt"
	_ce "github.com/AlexGames73/unioffice-free/common"
	_eb "github.com/AlexGames73/unioffice-free/common/logger"
	_d "io"
	_fa "io/ioutil"
	_cbc "log"
	_ee "net"
	_cc "net/http"
	_cd "os"
	_ed "path/filepath"
	_a "regexp"
	_f "sort"
	_ab "strings"
	_e "sync"
	_gg "time"
)

const _bdd = "\u0055\u004e\u0049OF\u0046\u0049\u0043\u0045\u005f\u0043\u0055\u0053\u0054\u004f\u004d\u0045\u0052\u005f\u004e\u0041\u004d\u0045"

type defaultStateHolder struct{}

func GetLicenseKey() *LicenseKey {
	if _bgab == nil {
		return nil
	}
	_cee := *_bgab
	return &_cee
}

type meteredUsageCheckinForm struct {
	Instance       string         `json:"inst"`
	Next           string         `json:"next"`
	UsageNumber    int            `json:"usage_number"`
	NumFailed      int64          `json:"num_failed"`
	Hostname       string         `json:"hostname"`
	LocalIP        string         `json:"local_ip"`
	MacAddress     string         `json:"mac_address"`
	Package        string         `json:"package"`
	PackageVersion string         `json:"package_version"`
	Usage          map[string]int `json:"u"`
}

func _cge(_agb, _eda []byte) ([]byte, error) {
	_dfa, _gae := _fe.NewCipher(_agb)
	if _gae != nil {
		return nil, _gae
	}
	_bdc := make([]byte, _fe.BlockSize+len(_eda))
	_bdb := _bdc[:_fe.BlockSize]
	if _, _aeb := _d.ReadFull(_ef.Reader, _bdb); _aeb != nil {
		return nil, _aeb
	}
	_fgde := _agf.NewCFBEncrypter(_dfa, _bdb)
	_fgde.XORKeyStream(_bdc[_fe.BlockSize:], _eda)
	_eebe := make([]byte, _bd.URLEncoding.EncodedLen(len(_bdc)))
	_bd.URLEncoding.Encode(_eebe, _bdc)
	return _eebe, nil
}

var _ege stateLoader = defaultStateHolder{}

type LicenseKey struct {
	LicenseId    string   `json:"license_id"`
	CustomerId   string   `json:"customer_id"`
	CustomerName string   `json:"customer_name"`
	Tier         string   `json:"tier"`
	CreatedAt    _gg.Time `json:"-"`
	CreatedAtInt int64    `json:"created_at"`
	ExpiresAt    _gg.Time `json:"-"`
	ExpiresAtInt int64    `json:"expires_at"`
	CreatedBy    string   `json:"created_by"`
	CreatorName  string   `json:"creator_name"`
	CreatorEmail string   `json:"creator_email"`
	UniPDF       bool     `json:"unipdf"`
	UniOffice    bool     `json:"unioffice"`
	UniHTML      bool     `json:"unihtml"`
	Trial        bool     `json:"trial"`
	_cfd         bool
	_eeb         string
}

type reportState struct {
	Instance      string         `json:"inst"`
	Next          string         `json:"n"`
	Docs          int64          `json:"d"`
	NumErrors     int64          `json:"e"`
	LimitDocs     bool           `json:"ld"`
	RemainingDocs int64          `json:"rd"`
	LastReported  _gg.Time       `json:"lr"`
	LastWritten   _gg.Time       `json:"lw"`
	Usage         map[string]int `json:"u"`
}

func _adba() ([]string, []string, error) {
	_agg, _cgaf := _ee.Interfaces()
	if _cgaf != nil {
		return nil, nil, _cgaf
	}
	var _afca []string
	var _dcd []string
	for _, _cdg := range _agg {
		if _cdg.Flags&_ee.FlagUp == 0 || _b.Equal(_cdg.HardwareAddr, nil) {
			continue
		}
		_bgcf, _edc := _cdg.Addrs()
		if _edc != nil {
			return nil, nil, _edc
		}
		_gba := 0
		for _, _cgab := range _bgcf {
			var _gebf _ee.IP
			switch _gcc := _cgab.(type) {
			case *_ee.IPNet:
				_gebf = _gcc.IP
			case *_ee.IPAddr:
				_gebf = _gcc.IP
			}
			if _gebf.IsLoopback() {
				continue
			}
			if _gebf.To4() == nil {
				continue
			}
			_dcd = append(_dcd, _gebf.String())
			_gba++
		}
		_ega := _cdg.HardwareAddr.String()
		if _ega != "" && _gba > 0 {
			_afca = append(_afca, _ega)
		}
	}
	return _afca, _dcd, nil
}

func _ae(_feg string, _aef []byte) (string, error) {
	_dga, _ := _ca.Decode([]byte(_feg))
	if _dga == nil {
		return "", _ag.Errorf("\u0050\u0072\u0069\u0076\u004b\u0065\u0079\u0020\u0066a\u0069\u006c\u0065\u0064")
	}
	_agcg, _fea := _fed.ParsePKCS1PrivateKey(_dga.Bytes)
	if _fea != nil {
		return "", _fea
	}
	_be := _abb.New()
	_be.Write(_aef)
	_bb := _be.Sum(nil)
	_faf, _fea := _df.SignPKCS1v15(_ef.Reader, _agcg, _g.SHA512, _bb)
	if _fea != nil {
		return "", _fea
	}
	_efb := _bd.StdEncoding.EncodeToString(_aef)
	_efb += "\u000a\u002b\u000a"
	_efb += _bd.StdEncoding.EncodeToString(_faf)
	return _efb, nil
}

const (
	_ff  = "\u002d\u002d\u002d--\u0042\u0045\u0047\u0049\u004e\u0020\u0055\u004e\u0049D\u004fC\u0020L\u0049C\u0045\u004e\u0053\u0045\u0020\u004b\u0045\u0059\u002d\u002d\u002d\u002d\u002d"
	_gga = "\u002d\u002d\u002d\u002d\u002d\u0045\u004e\u0044\u0020\u0055\u004e\u0049\u0044\u004f\u0043 \u004cI\u0043\u0045\u004e\u0053\u0045\u0020\u004b\u0045\u0059\u002d\u002d\u002d\u002d\u002d"
)

type meteredStatusForm struct{}

const _gge = "\u0033\u0030\u0035\u0063\u0033\u0030\u0030\u00640\u0036\u0030\u0039\u0032\u0061\u0038\u00364\u0038\u0038\u0036\u0066\u0037\u0030d\u0030\u0031\u0030\u0031\u0030\u00310\u0035\u0030\u0030\u0030\u0033\u0034\u0062\u0030\u0030\u0033\u0030\u00348\u0030\u0032\u0034\u0031\u0030\u0030\u0062\u0038\u0037\u0065\u0061\u0066\u0062\u0036\u0063\u0030\u0037\u0034\u0039\u0039\u0065\u0062\u00397\u0063\u0063\u0039\u0064\u0033\u0035\u0036\u0035\u0065\u0063\u00663\u0031\u0036\u0038\u0031\u0039\u0036\u0033\u0030\u0031\u0039\u0030\u0037c\u0038\u0034\u0031\u0061\u0064\u0064c6\u0036\u0035\u0030\u0038\u0036\u0062\u0062\u0033\u0065\u0064\u0038\u0065\u0062\u0031\u0032\u0064\u0039\u0064\u0061\u0032\u0036\u0063\u0061\u0066\u0061\u0039\u0036\u00345\u0030\u00314\u0036\u0064\u0061\u0038\u0062\u0064\u0030\u0063c\u0066\u0031\u0035\u0035\u0066\u0063a\u0063\u0063\u00368\u0036\u0039\u0035\u0035\u0065\u0066\u0030\u0033\u0030\u0032\u0066\u0061\u0034\u0034\u0061\u0061\u0033\u0065\u0063\u0038\u0039\u0034\u0031\u0037\u0062\u0030\u0032\u0030\u0033\u0030\u0031\u0030\u0030\u0030\u0031"

var _edf = false

func _gbc(_abgb *_cc.Response) ([]byte, error) {
	var _adaa []byte
	_afbc, _cgb := _cebe(_abgb)
	if _cgb != nil {
		return _adaa, _cgb
	}
	return _fa.ReadAll(_afbc)
}

func SetLegacyLicenseKey(s string) error {
	_gff := _a.MustCompile("\u005c\u0073")
	s = _gff.ReplaceAllString(s, "")
	var _daa _d.Reader
	_daa = _ab.NewReader(s)
	_daa = _bd.NewDecoder(_bd.RawURLEncoding, _daa)
	_daa,
		_ebe := _cb.NewReader(_daa)
	if _ebe != nil {
		return _ebe
	}
	_fee := _dg.NewDecoder(_daa)
	_faga := &LegacyLicense{}
	if _gcd := _fee.Decode(_faga); _gcd != nil {
		return _gcd
	}
	if _bbc := _faga.Verify(_aa); _bbc != nil {
		return _gb.New("\u006c\u0069\u0063en\u0073\u0065\u0020\u0076\u0061\u006c\u0069\u0064\u0061\u0074\u0069\u006e\u0020\u0065\u0072\u0072\u006f\u0072")
	}
	if _faga.Expiration.Before(_ce.ReleasedAt) {
		return _gb.New("\u006ci\u0063e\u006e\u0073\u0065\u0020\u0065\u0078\u0070\u0069\u0072\u0065\u0064")
	}
	_gee := _gg.Now().UTC()
	_abeg := LicenseKey{}
	_abeg.CreatedAt = _gee
	_abeg.CustomerId = "\u004c\u0065\u0067\u0061\u0063\u0079"
	_abeg.CustomerName = _faga.Name
	_abeg.Tier = LicenseTierBusiness
	_abeg.ExpiresAt = _faga.Expiration
	_abeg.CreatorName = "\u0055\u006e\u0069\u0044\u006f\u0063\u0020\u0073\u0075p\u0070\u006f\u0072\u0074"
	_abeg.CreatorEmail = "\u0073\u0075\u0070\u0070\u006f\u0072\u0074\u0040\u0075\u006e\u0069\u0064o\u0063\u002e\u0069\u006f"
	_abeg.UniOffice = true
	_bgab = &_abeg
	return nil
}

func _cebe(_acab *_cc.Response) (_d.ReadCloser, error) {
	var _edcd error
	var _egdf _d.ReadCloser
	switch _ab.ToLower(_acab.Header.Get("\u0043\u006fn\u0074\u0065\u006et\u002d\u0045\u006e\u0063\u006f\u0064\u0069\u006e\u0067")) {
	case "\u0067\u007a\u0069\u0070":
		_egdf, _edcd = _cb.NewReader(_acab.Body)
		if _edcd != nil {
			return _egdf, _edcd
		}
		defer _egdf.Close()
	default:
		_egdf = _acab.Body
	}
	return _egdf, nil
}

func SetLicenseKey(content string, customerName string) error {
	if _edf {
		return nil
	}
	_fddd, _fefa := _dedb(content)
	if _fefa != nil {
		_eb.Log.Error("\u004c\u0069c\u0065\u006e\u0073\u0065\u0020\u0063\u006f\u0064\u0065\u0020\u0064\u0065\u0063\u006f\u0064\u0065\u0020\u0065\u0072\u0072\u006f\u0072: \u0025\u0076", _fefa)
		return _fefa
	}
	if !_ab.EqualFold(_fddd.CustomerName, customerName) {
		_eb.Log.Error("L\u0069ce\u006es\u0065 \u0063\u006f\u0064\u0065\u0020i\u0073\u0073\u0075e\u0020\u002d\u0020\u0043\u0075s\u0074\u006f\u006de\u0072\u0020\u006e\u0061\u006d\u0065\u0020\u006d\u0069\u0073\u006da\u0074\u0063\u0068, e\u0078\u0070\u0065\u0063\u0074\u0065d\u0020\u0027\u0025\u0073\u0027\u002c\u0020\u0062\u0075\u0074\u0020\u0067o\u0074 \u0027\u0025\u0073\u0027", customerName, _fddd.CustomerName)
		return _ag.Errorf("\u0063\u0075\u0073\u0074\u006fm\u0065\u0072\u0020\u006e\u0061\u006d\u0065\u0020\u006d\u0069\u0073\u006d\u0061t\u0063\u0068\u002c\u0020\u0065\u0078\u0070\u0065\u0063\u0074\u0065\u0064\u0020\u0027\u0025\u0073\u0027\u002c\u0020\u0062\u0075\u0074\u0020\u0067\u006f\u0074\u0020\u0027\u0025\u0073'", customerName, _fddd.CustomerName)
	}
	_fefa = _fddd.Validate()
	if _fefa != nil {
		_eb.Log.Error("\u004c\u0069\u0063\u0065\u006e\u0073e\u0020\u0063\u006f\u0064\u0065\u0020\u0076\u0061\u006c\u0069\u0064\u0061\u0074i\u006f\u006e\u0020\u0065\u0072\u0072\u006fr\u003a\u0020\u0025\u0076", _fefa)
		return _fefa
	}
	_bgab = &_fddd
	return nil
}

func Track(docKey string, useKey string) error {
	return _fbgb(docKey, useKey, false)
}

func (_age *LicenseKey) isExpired() bool {
	return _age.getExpiryDateToCompare().After(_age.ExpiresAt)
}

func init() {
	_fca, _ea := _agc.DecodeString(_gge)
	if _ea != nil {
		_cbc.Fatalf("e\u0072\u0072\u006f\u0072 r\u0065a\u0064\u0069\u006e\u0067\u0020k\u0065\u0079\u003a\u0020\u0025\u0073", _ea)
	}
	_dfg, _ea := _fed.ParsePKIXPublicKey(_fca)
	if _ea != nil {
		_cbc.Fatalf("e\u0072\u0072\u006f\u0072 r\u0065a\u0064\u0069\u006e\u0067\u0020k\u0065\u0079\u003a\u0020\u0025\u0073", _ea)
	}
	_aa = _dfg.(*_df.PublicKey)
}

type meteredStatusResp struct {
	Valid        bool  `json:"valid"`
	OrgCredits   int64 `json:"org_credits"`
	OrgUsed      int64 `json:"org_used"`
	OrgRemaining int64 `json:"org_remaining"`
}

var _egd map[string]struct{}

var _eef = _gg.Date(2010, 1, 1, 0, 0, 0, 0, _gg.UTC)

func GenRefId(prefix string) (string, error) {
	var _fdb _b.Buffer
	_fdb.WriteString(prefix)
	_ebb := make([]byte, 8+16)
	_dbff := _gg.Now().UTC().UnixNano()
	_db.BigEndian.PutUint64(_ebb, uint64(_dbff))
	_, _dedbf := _ef.Read(_ebb[8:])
	if _dedbf != nil {
		return "", _dedbf
	}
	_fdb.WriteString(_agc.EncodeToString(_ebb))
	return _fdb.String(), nil
}

const _bbae = "\u0055\u004e\u0049\u004fFF\u0049\u0043\u0045\u005f\u004c\u0049\u0043\u0045\u004e\u0053\u0045\u005f\u0050\u0041T\u0048"

func SetMeteredKey(apiKey string) error {
	if len(apiKey) == 0 {
		_eb.Log.Error("\u004d\u0065\u0074\u0065\u0072e\u0064\u0020\u004c\u0069\u0063\u0065\u006e\u0073\u0065\u0020\u0041\u0050\u0049 \u004b\u0065\u0079\u0020\u006d\u0075\u0073\u0074\u0020\u006e\u006f\u0074\u0020\u0062\u0065\u0020\u0065\u006d\u0070\u0074\u0079")
		_eb.Log.Error("\u002d\u0020\u0047\u0072\u0061\u0062\u0020\u006f\u006e\u0065\u0020\u0069\u006e\u0020\u0074h\u0065\u0020\u0046\u0072\u0065\u0065\u0020\u0054\u0069\u0065\u0072\u0020\u0061t\u0020\u0068\u0074\u0074\u0070\u0073\u003a\u002f\u002f\u0063\u006c\u006fud\u002e\u0075\u006e\u0069\u0064\u006f\u0063\u002e\u0069\u006f")
		return _ag.Errorf("\u006de\u0074\u0065\u0072e\u0064\u0020\u006ci\u0063en\u0073\u0065\u0020\u0061\u0070\u0069\u0020k\u0065\u0079\u0020\u006d\u0075\u0073\u0074\u0020\u006e\u006f\u0074\u0020\u0062\u0065\u0020\u0065\u006d\u0070\u0074\u0079\u003a\u0020\u0063\u0072\u0065\u0061\u0074\u0065 o\u006ee\u0020\u0061\u0074\u0020\u0068\u0074t\u0070\u0073\u003a\u002f\u002fc\u006c\u006f\u0075\u0064\u002e\u0075\u006e\u0069\u0064\u006f\u0063.\u0069\u006f")
	}
	if _bgab != nil && (_bgab._cfd || _bgab.Tier != LicenseTierUnlicensed) {
		_eb.Log.Error("\u0045\u0052\u0052\u004f\u0052:\u0020\u0043\u0061\u006e\u006eo\u0074 \u0073\u0065\u0074\u0020\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u0020\u006b\u0065\u0079\u0020\u0074\u0077\u0069c\u0065\u0020\u002d\u0020\u0053\u0068\u006f\u0075\u006c\u0064\u0020\u006a\u0075\u0073\u0074\u0020\u0069\u006e\u0069\u0074\u0069\u0061\u006c\u0069z\u0065\u0020\u006f\u006e\u0063\u0065")
		return _gb.New("\u006c\u0069\u0063en\u0073\u0065\u0020\u006b\u0065\u0079\u0020\u0061\u006c\u0072\u0065\u0061\u0064\u0079\u0020\u0073\u0065\u0074")
	}
	_dfd := _cgf()
	_dfd._dac = apiKey
	_bae, _edd := _dfd.getStatus()
	if _edd != nil {
		return _edd
	}
	if !_bae.Valid {
		return _gb.New("\u006b\u0065\u0079\u0020\u006e\u006f\u0074\u0020\u0076\u0061\u006c\u0069\u0064")
	}
	_agef := &LicenseKey{
		_cfd: true,
		_eeb: apiKey,
	}
	_bgab = _agef
	return nil
}

const _ffeg = "\u000a\u002d\u002d\u002d\u002d\u002d\u0042\u0045\u0047\u0049\u004e \u0050\u0055\u0042\u004c\u0049\u0043\u0020\u004b\u0045Y\u002d\u002d\u002d\u002d\u002d\u000a\u004d\u0049I\u0042\u0049\u006a\u0041NB\u0067\u006b\u0071\u0068\u006b\u0069G\u0039\u0077\u0030\u0042\u0041\u0051\u0045\u0046A\u0041\u004f\u0043\u0041\u0051\u0038\u0041\u004d\u0049\u0049\u0042\u0043\u0067\u004b\u0043\u0041\u0051\u0045A\u006dF\u0055\u0069\u0079\u0064\u0037\u0062\u0035\u0058\u006a\u0070\u006b\u0050\u0035\u0052\u0061\u0070\u0034\u0077\u000a\u0044\u0063\u0031d\u0079\u007a\u0049\u0051\u0034\u004c\u0065\u006b\u0078\u0072\u0076\u0079\u0074\u006e\u0045\u004d\u0070\u004e\u0055\u0062\u006f\u0036i\u0041\u0037\u0034\u0056\u0038\u0072\u0075\u005a\u004f\u0076\u0072\u0053\u0063\u0073\u0066\u0032\u0051\u0065\u004e9\u002f\u0071r\u0055\u0047\u0038\u0071\u0045\u0062\u0055\u0057\u0064\u006f\u0045\u0059\u0071+\u000a\u006f\u0074\u0046\u004e\u0041\u0046N\u0078\u006c\u0047\u0062\u0078\u0062\u0044\u0048\u0063\u0064\u0047\u0056\u0061\u004d\u0030\u004f\u0058\u0064\u0058g\u0044y\u004c5\u0061\u0049\u0045\u0061\u0067\u004c\u0030\u0063\u0035\u0070\u0077\u006a\u0049\u0064\u0050G\u0049\u006e\u0034\u0036\u0066\u0037\u0038\u0065\u004d\u004a\u002b\u004a\u006b\u0064\u0063\u0070\u0044\n\u0044\u004a\u0061\u0071\u0059\u0058d\u0072\u007a5\u004b\u0065\u0073\u0068\u006aS\u0069\u0049\u0061\u0061\u0037\u006d\u0065\u006e\u0042\u0049\u0041\u0058\u0053\u0034\u0055\u0046\u0078N\u0066H\u0068\u004e\u0030\u0048\u0043\u0059\u005a\u0059\u0071\u0051\u0047\u0037\u0062K+\u0073\u0035\u0072R\u0048\u006f\u006e\u0079\u0064\u004eW\u0045\u0047\u000a\u0048\u0038M\u0079\u0076\u00722\u0070\u0079\u0061\u0032K\u0072\u004d\u0075m\u0066\u006d\u0041\u0078\u0055\u0042\u0036\u0066\u0065\u006e\u0043\u002f4\u004f\u0030\u0057\u00728\u0067\u0066\u0050\u004f\u0055\u0038R\u0069\u0074\u006d\u0062\u0044\u0076\u0051\u0050\u0049\u0052\u0058\u004fL\u0034\u0076\u0054B\u0072\u0042\u0064\u0062a\u0041\u000a9\u006e\u0077\u004e\u0050\u002b\u0069\u002f\u002f\u0032\u0030\u004d\u00542\u0062\u0078\u006d\u0065\u0057\u0042\u002b\u0067\u0070\u0063\u0045\u0068G\u0070\u0058\u005a7\u0033\u0033\u0061\u007a\u0051\u0078\u0072\u0043\u0033\u004a\u0034\u0076\u0033C\u005a\u006d\u0045\u004eS\u0074\u0044\u004b\u002f\u004b\u0044\u0053\u0050\u004b\u0055\u0047\u0066\u00756\u000a\u0066\u0077I\u0044\u0041\u0051\u0041\u0042\u000a\u002d\u002d\u002d\u002d\u002dE\u004e\u0044\u0020\u0050\u0055\u0042\u004c\u0049\u0043 \u004b\u0045Y\u002d\u002d\u002d\u002d\u002d\n"

type stateLoader interface {
	loadState(_fede string) (reportState, error)
	updateState(_cbf, _efef, _aegd string, _cccg int, _ebcg bool, _ecdg int, _dea int, _aaa _gg.Time, _ddc map[string]int) error
}

var _bgab = MakeUnlicensedKey()

type LegacyLicenseType byte

func _eae(_cba, _dba []byte) ([]byte, error) {
	_bcd := make([]byte, _bd.URLEncoding.DecodedLen(len(_dba)))
	_aaf,
		_gbae := _bd.URLEncoding.Decode(_bcd, _dba)
	if _gbae != nil {
		return nil, _gbae
	}
	_bcd = _bcd[:_aaf]
	_ece,
		_gbae := _fe.NewCipher(_cba)
	if _gbae != nil {
		return nil, _gbae
	}
	if len(_bcd) < _fe.BlockSize {
		return nil, _gb.New("c\u0069p\u0068\u0065\u0072\u0074\u0065\u0078\u0074\u0020t\u006f\u006f\u0020\u0073ho\u0072\u0074")
	}
	_dfb := _bcd[:_fe.BlockSize]
	_bcd = _bcd[_fe.BlockSize:]
	_bag := _agf.NewCFBDecrypter(_ece, _dfb)
	_bag.XORKeyStream(_bcd, _bcd)
	return _bcd,
		nil
}

var _gde = &_e.Mutex{}

type LegacyLicense struct {
	Name        string
	Signature   string `json:",omitempty"`
	Expiration  _gg.Time
	LicenseType LegacyLicenseType
}

func _fec(_eeed []byte) (_d.Reader, error) {
	_caa := new(_b.Buffer)
	_eefdc := _cb.NewWriter(_caa)
	_eefdc.Write(_eeed)
	_dafa := _eefdc.Close()
	if _dafa != nil {
		return nil, _dafa
	}
	return _caa,
		nil
}

func (_ec *LicenseKey) Validate() error {
	if _ec._cfd {
		return nil
	}
	if len(_ec.LicenseId) < 10 {
		return _ag.Errorf("i\u006e\u0076\u0061\u006c\u0069\u0064 \u006c\u0069\u0063\u0065\u006e\u0073\u0065\u003a\u0020L\u0069\u0063\u0065n\u0073e\u0020\u0049\u0064")
	}
	if len(_ec.CustomerId) < 10 {
		return _ag.Errorf("\u0069\u006e\u0076\u0061l\u0069\u0064\u0020\u006c\u0069\u0063\u0065\u006e\u0073\u0065:\u0020C\u0075\u0073\u0074\u006f\u006d\u0065\u0072 \u0049\u0064")
	}
	if len(_ec.CustomerName) < 1 {
		return _ag.Errorf("\u0069\u006e\u0076\u0061\u006c\u0069\u0064\u0020\u006c\u0069c\u0065\u006e\u0073\u0065\u003a\u0020\u0043u\u0073\u0074\u006f\u006d\u0065\u0072\u0020\u004e\u0061\u006d\u0065")
	}
	if _eef.After(_ec.CreatedAt) {
		return _ag.Errorf("\u0069\u006e\u0076\u0061\u006c\u0069\u0064\u0020\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u003a\u0020\u0043\u0072\u0065\u0061\u0074\u0065\u0064 \u0041\u0074\u0020\u0069\u0073 \u0069\u006ev\u0061\u006c\u0069\u0064")
	}
	if _ec.ExpiresAt.IsZero() {
		_ad := _ec.CreatedAt.AddDate(1, 0, 0)
		if _efg.After(_ad) {
			_ad = _efg
		}
		_ec.ExpiresAt = _ad
	}
	if _ec.CreatedAt.After(_ec.ExpiresAt) {
		return _ag.Errorf("i\u006e\u0076\u0061\u006c\u0069\u0064\u0020\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u003a\u0020\u0043\u0072\u0065\u0061\u0074\u0065\u0064\u0020\u0041\u0074 \u0063a\u006e\u006e\u006f\u0074 \u0062\u0065 \u0047\u0072\u0065\u0061\u0074\u0065\u0072\u0020\u0074\u0068\u0061\u006e\u0020\u0045\u0078\u0070\u0069\u0072\u0065\u0073\u0020\u0041\u0074")
	}
	if _ec.isExpired() {
		return _ag.Errorf("\u0069\u006e\u0076\u0061\u006c\u0069\u0064\u0020l\u0069\u0063\u0065ns\u0065\u003a\u0020\u0054\u0068\u0065 \u006c\u0069\u0063\u0065\u006e\u0073\u0065\u0020\u0068\u0061\u0073\u0020\u0061\u006c\u0072e\u0061\u0064\u0079\u0020\u0065\u0078\u0070\u0069r\u0065\u0064")
	}
	if len(_ec.CreatorName) < 1 {
		return _ag.Errorf("\u0069\u006ev\u0061\u006c\u0069\u0064\u0020\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u003a\u0020\u0043\u0072\u0065\u0061\u0074\u006f\u0072\u0020na\u006d\u0065")
	}
	if len(_ec.CreatorEmail) < 1 {
		return _ag.Errorf("\u0069\u006e\u0076\u0061\u006c\u0069\u0064\u0020\u006c\u0069c\u0065\u006e\u0073\u0065\u003a\u0020\u0043r\u0065\u0061\u0074\u006f\u0072\u0020\u0065\u006d\u0061\u0069\u006c")
	}
	if _ec.CreatedAt.After(_fd) {
		if !_ec.UniOffice {
			return _ag.Errorf("\u0069\u006e\u0076\u0061l\u0069\u0064\u0020\u006c\u0069\u0063\u0065\u006e\u0073e\u003a\u0020\u0054\u0068\u0069\u0073\u0020\u0055\u006e\u0069\u0044\u006f\u0063\u0020\u006b\u0065\u0079\u0020i\u0073\u0020\u0069\u006e\u0076a\u006c\u0069\u0064\u0020\u0066\u006f\u0072\u0020\u0055\u006e\u0069\u004f\u0066\u0066\u0069\u0063\u0065")
		}
	}
	return nil
}

func _fbgb(_gcg string, _aad string, _egeg bool) error {
	if _bgab == nil {
		return _gb.New("\u006e\u006f\u0020\u006c\u0069\u0063\u0065\u006e\u0073e\u0020\u006b\u0065\u0079")
	}
	if !_bgab._cfd || len(_bgab._eeb) == 0 {
		return nil
	}
	if len(_gcg) == 0 && !_egeg {
		return _gb.New("\u0064\u006f\u0063\u004b\u0065\u0079\u0020\u006e\u006ft\u0020\u0073\u0065\u0074")
	}
	_gde.Lock()
	defer _gde.Unlock()
	if _egd == nil {
		_egd = map[string]struct{}{}
	}
	if _bgb == nil {
		_bgb = map[string]int{}
	}
	_fbf := 0
	if !_egeg {
		_, _gfc := _egd[_gcg]
		if !_gfc {
			_egd[_gcg] = struct{}{}
			_fbf++
		}
		if _fbf == 0 {
			return nil
		}
		_bgb[_aad]++
	}
	_cda := _gg.Now()
	_ffe, _geb := _ege.loadState(_bgab._eeb)
	if _geb != nil {
		_eb.Log.Error("\u0045R\u0052\u004f\u0052\u003a\u0020\u0025v", _geb)
		return _geb
	}
	if _ffe.Usage == nil {
		_ffe.Usage = map[string]int{}
	}
	for _beaf, _eefd := range _bgb {
		_ffe.Usage[_beaf] += _eefd
	}
	_bgb = nil
	const _dge = 24 * _gg.Hour
	const _dbg = 3 * 24 * _gg.Hour
	if len(_ffe.Instance) == 0 || _cda.Sub(_ffe.LastReported) > _dge || (_ffe.LimitDocs && _ffe.RemainingDocs <= _ffe.Docs+int64(_fbf)) || _egeg {
		_fafd, _fab := _cd.Hostname()
		if _fab != nil {
			return _fab
		}
		_bgc := _ffe.Docs
		_dae, _aae, _fab := _adba()
		if _fab != nil {
			return _fab
		}
		_f.Strings(_aae)
		_f.Strings(_dae)
		_egegf, _fab := _afce()
		if _fab != nil {
			return _fab
		}
		_eab := false
		for _, _bce := range _aae {
			if _bce == _egegf.String() {
				_eab = true
			}
		}
		if !_eab {
			_aae = append(_aae, _egegf.String())
		}
		_afbb := _cgf()
		_afbb._dac = _bgab._eeb
		_bgc += int64(_fbf)
		_dfdc := meteredUsageCheckinForm{
			Instance:       _ffe.Instance,
			Next:           _ffe.Next,
			UsageNumber:    int(_bgc),
			NumFailed:      _ffe.NumErrors,
			Hostname:       _fafd,
			LocalIP:        _ab.Join(_aae, "\u002c\u0020"),
			MacAddress:     _ab.Join(_dae, "\u002c\u0020"),
			Package:        "\u0075n\u0069\u006f\u0066\u0066\u0069\u0063e",
			PackageVersion: _ce.Version,
			Usage:          _ffe.Usage,
		}
		if len(_dae) == 0 {
			_dfdc.MacAddress = "\u006e\u006f\u006e\u0065"
		}
		_bfaf := int64(0)
		_bbd := _ffe.NumErrors
		_ffb := _cda
		_geg := 0
		_ged := _ffe.LimitDocs
		_ecb, _fab := _afbb.checkinUsage(_dfdc)
		if _fab != nil {
			if _cda.Sub(_ffe.LastReported) > _dbg {
				return _gb.New("\u0074\u006f\u006f\u0020\u006c\u006f\u006e\u0067\u0020\u0073\u0069\u006e\u0063\u0065\u0020\u006c\u0061\u0073\u0074\u0020\u0073\u0075\u0063\u0063e\u0073\u0073\u0066\u0075\u006c \u0063\u0068e\u0063\u006b\u0069\u006e")
			}
			_bfaf = _bgc
			_bbd++
			_ffb = _ffe.LastReported
		} else {
			_ged = _ecb.LimitDocs
			_geg = _ecb.RemainingDocs
			_bbd = 0
		}
		if len(_ecb.Instance) == 0 {
			_ecb.Instance = _dfdc.Instance
		}
		if len(_ecb.Next) == 0 {
			_ecb.Next = _dfdc.Next
		}
		_fab = _ege.updateState(_afbb._dac, _ecb.Instance, _ecb.Next, int(_bfaf), _ged, _geg, int(_bbd), _ffb, nil)
		if _fab != nil {
			return _fab
		}
		if !_ecb.Success {
			return _ag.Errorf("\u0065r\u0072\u006f\u0072\u003a\u0020\u0025s", _ecb.Message)
		}
	} else {
		_geb = _ege.updateState(_bgab._eeb, _ffe.Instance, _ffe.Next, int(_ffe.Docs)+_fbf, _ffe.LimitDocs, int(_ffe.RemainingDocs), int(_ffe.NumErrors), _ffe.LastReported, _ffe.Usage)
		if _geb != nil {
			return _geb
		}
	}
	return nil
}

func (_bg *meteredClient) getStatus() (meteredStatusResp, error) {
	var _cac meteredStatusResp
	_afg := _bg._fbg + "\u002fm\u0065t\u0065\u0072\u0065\u0064\u002f\u0073\u0074\u0061\u0074\u0075\u0073"
	var _bfg meteredStatusForm
	_ecd, _cab := _dg.Marshal(_bfg)
	if _cab != nil {
		return _cac, _cab
	}
	_gbf, _cab := _fec(_ecd)
	if _cab != nil {
		return _cac, _cab
	}
	_dcf, _cab := _cc.NewRequest("\u0050\u004f\u0053\u0054", _afg, _gbf)
	if _cab != nil {
		return _cac, _cab
	}
	_dcf.Header.Add("\u0043\u006f\u006et\u0065\u006e\u0074\u002d\u0054\u0079\u0070\u0065", "\u0061\u0070p\u006c\u0069\u0063a\u0074\u0069\u006f\u006e\u002f\u006a\u0073\u006f\u006e")
	_dcf.Header.Add("\u0043\u006fn\u0074\u0065\u006et\u002d\u0045\u006e\u0063\u006f\u0064\u0069\u006e\u0067", "\u0067\u007a\u0069\u0070")
	_dcf.Header.Add("\u0041c\u0063e\u0070\u0074\u002d\u0045\u006e\u0063\u006f\u0064\u0069\u006e\u0067", "\u0067\u007a\u0069\u0070")
	_dcf.Header.Add("\u0058-\u0041\u0050\u0049\u002d\u004b\u0045Y", _bg._dac)
	_ggc, _cab := _bg._ccd.Do(_dcf)
	if _cab != nil {
		return _cac, _cab
	}
	defer _ggc.Body.Close()
	if _ggc.StatusCode != 200 {
		return _cac, _ag.Errorf("\u0066\u0061i\u006c\u0065\u0064\u0020t\u006f\u0020c\u0068\u0065\u0063\u006b\u0069\u006e\u002c\u0020s\u0074\u0061\u0074\u0075\u0073\u0020\u0063\u006f\u0064\u0065\u0020\u0069s\u003a\u0020\u0025\u0064", _ggc.StatusCode)
	}
	_ac, _cab := _gbc(_ggc)
	if _cab != nil {
		return _cac, _cab
	}
	_cab = _dg.Unmarshal(_ac, &_cac)
	if _cab != nil {
		return _cac, _cab
	}
	return _cac, nil
}

var _fd = _gg.Date(2019, 6, 6, 0, 0, 0, 0, _gg.UTC)

func (_fg LegacyLicense) Verify(pubKey *_df.PublicKey) error {
	_fb := _fg
	_fb.Signature = ""
	_eea := _b.Buffer{}
	_bed := _dg.NewEncoder(&_eea)
	if _bbb := _bed.Encode(_fb); _bbb != nil {
		return _bbb
	}
	_eee,
		_cfg := _agc.DecodeString(_fg.Signature)
	if _cfg != nil {
		return _cfg
	}
	_ceb := _ge.Sum256(_eea.Bytes())
	_cfg = _df.VerifyPKCS1v15(pubKey, _g.SHA256, _ceb[:], _eee)
	return _cfg
}

func _dedb(_gf string) (LicenseKey, error) {
	var _eg LicenseKey
	_fcd, _dbf := _bad(_ff, _gga, _gf)
	if _dbf != nil {
		return _eg, _dbf
	}
	_ebc, _dbf := _efe(_ffeg, _fcd)
	if _dbf != nil {
		return _eg, _dbf
	}
	_dbf = _dg.Unmarshal(_ebc, &_eg)
	if _dbf != nil {
		return _eg, _dbf
	}
	_eg.CreatedAt = _gg.Unix(_eg.CreatedAtInt, 0)
	if _eg.ExpiresAtInt > 0 {
		_edb := _gg.Unix(_eg.ExpiresAtInt, 0)
		_eg.ExpiresAt = _edb
	}
	return _eg, nil
}

func TrackUse(useKey string) {
	if _bgab == nil {
		return
	}
	if !_bgab._cfd || len(_bgab._eeb) == 0 {
		return
	}
	if len(useKey) == 0 {
		return
	}
	_gde.Lock()
	defer _gde.Unlock()
	if _bgb == nil {
		_bgb = map[string]int{}
	}
	_bgb[useKey]++
}

type meteredClient struct {
	_fbg string
	_dac string
	_ccd *_cc.Client
}

type meteredUsageCheckinResp struct {
	Instance      string `json:"inst"`
	Next          string `json:"next"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	RemainingDocs int    `json:"rd"`
	LimitDocs     bool   `json:"ld"`
}

func GetMeteredState() (MeteredStatus, error) {
	if _bgab == nil {
		return MeteredStatus{}, _gb.New("\u006c\u0069\u0063\u0065ns\u0065\u0020\u006b\u0065\u0079\u0020\u006e\u006f\u0074\u0020\u0073\u0065\u0074")
	}
	if !_bgab._cfd || len(_bgab._eeb) == 0 {
		return MeteredStatus{}, _gb.New("\u0061p\u0069 \u006b\u0065\u0079\u0020\u006e\u006f\u0074\u0020\u0073\u0065\u0074")
	}
	_aeg, _bedc := _ege.loadState(_bgab._eeb)
	if _bedc != nil {
		_eb.Log.Error("\u0045R\u0052\u004f\u0052\u003a\u0020\u0025v", _bedc)
		return MeteredStatus{}, _bedc
	}
	if _aeg.Docs > 0 {
		_afb := _fbgb("", "", true)
		if _afb != nil {
			return MeteredStatus{}, _afb
		}
	}
	_gde.Lock()
	defer _gde.Unlock()
	_dace := _cgf()
	_dace._dac = _bgab._eeb
	_fdd, _bedc := _dace.getStatus()
	if _bedc != nil {
		return MeteredStatus{}, _bedc
	}
	if !_fdd.Valid {
		return MeteredStatus{}, _gb.New("\u006b\u0065\u0079\u0020\u006e\u006f\u0074\u0020\u0076\u0061\u006c\u0069\u0064")
	}
	_efgd := MeteredStatus{
		OK:      true,
		Credits: _fdd.OrgCredits,
		Used:    _fdd.OrgUsed,
	}
	return _efgd, nil
}

var _aa *_df.PublicKey

func (_cae *LicenseKey) IsLicensed() bool {
	return true
}

func _cgf() *meteredClient {
	_gec := meteredClient{
		_fbg: "h\u0074\u0074\u0070\u0073\u003a\u002f/\u0063\u006c\u006f\u0075\u0064\u002e\u0075\u006e\u0069d\u006f\u0063\u002ei\u006f/\u0061\u0070\u0069",
		_ccd: &_cc.Client{
			Timeout: 30 * _gg.Second,
		},
	}
	if _gddg := _cd.Getenv("\u0055N\u0049\u0044\u004f\u0043_\u004c\u0049\u0043\u0045\u004eS\u0045_\u0053E\u0052\u0056\u0045\u0052\u005f\u0055\u0052L"); _ab.HasPrefix(_gddg, "\u0068\u0074\u0074\u0070") {
		_gec._fbg = _gddg
	}
	return &_gec
}

type MeteredStatus struct {
	OK      bool
	Credits int64
	Used    int64
}

func (_dgf *meteredClient) checkinUsage(_ccc meteredUsageCheckinForm) (meteredUsageCheckinResp, error) {
	_ccc.Package = "\u0075n\u0069\u006f\u0066\u0066\u0069\u0063e"
	_ccc.PackageVersion = _ce.Version
	var _cff meteredUsageCheckinResp
	_gc := _dgf._fbg + "\u002f\u006d\u0065\u0074er\u0065\u0064\u002f\u0075\u0073\u0061\u0067\u0065\u005f\u0063\u0068\u0065\u0063\u006bi\u006e"
	_ecg, _dfga := _dg.Marshal(_ccc)
	if _dfga != nil {
		return _cff, _dfga
	}
	_cga, _dfga := _fec(_ecg)
	if _dfga != nil {
		return _cff, _dfga
	}
	_bgd, _dfga := _cc.NewRequest("\u0050\u004f\u0053\u0054", _gc, _cga)
	if _dfga != nil {
		return _cff, _dfga
	}
	_bgd.Header.Add("\u0043\u006f\u006et\u0065\u006e\u0074\u002d\u0054\u0079\u0070\u0065", "\u0061\u0070p\u006c\u0069\u0063a\u0074\u0069\u006f\u006e\u002f\u006a\u0073\u006f\u006e")
	_bgd.Header.Add("\u0043\u006fn\u0074\u0065\u006et\u002d\u0045\u006e\u0063\u006f\u0064\u0069\u006e\u0067", "\u0067\u007a\u0069\u0070")
	_bgd.Header.Add("\u0041c\u0063e\u0070\u0074\u002d\u0045\u006e\u0063\u006f\u0064\u0069\u006e\u0067", "\u0067\u007a\u0069\u0070")
	_bgd.Header.Add("\u0058-\u0041\u0050\u0049\u002d\u004b\u0045Y", _dgf._dac)
	_aca, _dfga := _dgf._ccd.Do(_bgd)
	if _dfga != nil {
		return _cff, _dfga
	}
	defer _aca.Body.Close()
	if _aca.StatusCode != 200 {
		return _cff, _ag.Errorf("\u0066\u0061i\u006c\u0065\u0064\u0020t\u006f\u0020c\u0068\u0065\u0063\u006b\u0069\u006e\u002c\u0020s\u0074\u0061\u0074\u0075\u0073\u0020\u0063\u006f\u0064\u0065\u0020\u0069s\u003a\u0020\u0025\u0064", _aca.StatusCode)
	}
	_dda, _dfga := _gbc(_aca)
	if _dfga != nil {
		return _cff, _dfga
	}
	_dfga = _dg.Unmarshal(_dda, &_cff)
	if _dfga != nil {
		return _cff, _dfga
	}
	return _cff, nil
}

func _afce() (_ee.IP, error) {
	_daf, _eac := _ee.Dial("\u0075\u0064\u0070", "\u0038\u002e\u0038\u002e\u0038\u002e\u0038\u003a\u0038\u0030")
	if _eac != nil {
		return nil, _eac
	}
	defer _daf.Close()
	_ada := _daf.LocalAddr().(*_ee.UDPAddr)
	return _ada.IP, nil
}

func _aegb() string {
	_cea := _cd.Getenv("\u0048\u004f\u004d\u0045")
	if len(_cea) == 0 {
		_cea, _ = _cd.UserHomeDir()
	}
	return _cea
}

func _efe(_gbe string, _fc string) ([]byte, error) {
	var (
		_fcc int
		_dff string
	)
	for _, _dff = range []string{
		"\u000a\u002b\u000a", "\u000d\u000a\u002b\r\u000a", "\u0020\u002b\u0020",
	} {
		if _fcc = _ab.Index(_fc, _dff); _fcc != -1 {
			break
		}
	}
	if _fcc == -1 {
		return nil, _ag.Errorf("\u0069\u006e\u0076al\u0069\u0064\u0020\u0069\u006e\u0070\u0075\u0074\u002c \u0073i\u0067n\u0061t\u0075\u0072\u0065\u0020\u0073\u0065\u0070\u0061\u0072\u0061\u0074\u006f\u0072")
	}
	_cg := _fc[:_fcc]
	_bf := _fcc + len(_dff)
	_baa := _fc[_bf:]
	if _cg == "" || _baa == "" {
		return nil, _ag.Errorf("\u0069n\u0076\u0061l\u0069\u0064\u0020\u0069n\u0070\u0075\u0074,\u0020\u006d\u0069\u0073\u0073\u0069\u006e\u0067\u0020or\u0069\u0067\u0069n\u0061\u006c \u006f\u0072\u0020\u0073\u0069\u0067n\u0061\u0074u\u0072\u0065")
	}
	_ccf, _fag := _bd.StdEncoding.DecodeString(_cg)
	if _fag != nil {
		return nil, _ag.Errorf("\u0069\u006e\u0076\u0061li\u0064\u0020\u0069\u006e\u0070\u0075\u0074\u0020\u006f\u0072\u0069\u0067\u0069\u006ea\u006c")
	}
	_de, _fag := _bd.StdEncoding.DecodeString(_baa)
	if _fag != nil {
		return nil, _ag.Errorf("\u0069\u006e\u0076al\u0069\u0064\u0020\u0069\u006e\u0070\u0075\u0074\u0020\u0073\u0069\u0067\u006e\u0061\u0074\u0075\u0072\u0065")
	}
	_af, _ := _ca.Decode([]byte(_gbe))
	if _af == nil {
		return nil, _ag.Errorf("\u0050\u0075\u0062\u004b\u0065\u0079\u0020\u0066\u0061\u0069\u006c\u0065\u0064")
	}
	_cf, _fag := _fed.ParsePKIXPublicKey(_af.Bytes)
	if _fag != nil {
		return nil, _fag
	}
	_fagg := _cf.(*_df.PublicKey)
	if _fagg == nil {
		return nil, _ag.Errorf("\u0050u\u0062\u004b\u0065\u0079\u0020\u0063\u006f\u006e\u0076\u0065\u0072s\u0069\u006f\u006e\u0020\u0066\u0061\u0069\u006c\u0065\u0064")
	}
	_abe := _abb.New()
	_abe.Write(_ccf)
	_eba := _abe.Sum(nil)
	_fag = _df.VerifyPKCS1v15(_fagg, _g.SHA512, _eba, _de)
	if _fag != nil {
		return nil, _fag
	}
	return _ccf, nil
}

var _efg = _gg.Date(2020, 1, 1, 0, 0, 0, 0, _gg.UTC)

func init() {
	_acgg := _cd.Getenv(_bbae)
	_caba := _cd.Getenv(_bdd)
	if len(_acgg) == 0 || len(_caba) == 0 {
		return
	}
	_bedd,
		_egea := _fa.ReadFile(_acgg)
	if _egea != nil {
		_eb.Log.Error("\u0055\u006eab\u006c\u0065\u0020t\u006f\u0020\u0072\u0065ad \u006cic\u0065\u006e\u0073\u0065\u0020\u0063\u006fde\u0020\u0066\u0069\u006c\u0065\u003a\u0020%\u0076", _egea)
		return
	}
	_egea = SetLicenseKey(string(_bedd), _caba)
	if _egea != nil {
		_eb.Log.Error("\u0055\u006e\u0061b\u006c\u0065\u0020\u0074o\u0020\u006c\u006f\u0061\u0064\u0020\u006ci\u0063\u0065\u006e\u0073\u0065\u0020\u0063\u006f\u0064\u0065\u003a\u0020\u0025\u0076", _egea)
		return
	}
}

func MakeUnlicensedKey() *LicenseKey {
	_bea := LicenseKey{}
	_bea.CustomerName = "\u0055\u006e\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u0064"
	_bea.Tier = LicenseTierUnlicensed
	_bea.CreatedAt = _gg.Now().UTC()
	_bea.CreatedAtInt = _bea.CreatedAt.Unix()
	return &_bea
}

const (
	LicenseTierUnlicensed = "\u0075\u006e\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u0064"
	LicenseTierCommunity  = "\u0063o\u006d\u006d\u0075\u006e\u0069\u0074y"
	LicenseTierIndividual = "\u0069\u006e\u0064\u0069\u0076\u0069\u0064\u0075\u0061\u006c"
	LicenseTierBusiness   = "\u0062\u0075\u0073\u0069\u006e\u0065\u0073\u0073"
)

func (_badg *LicenseKey) ToString() string {
	if _badg._cfd {
		return "M\u0065t\u0065\u0072\u0065\u0064\u0020\u0073\u0075\u0062s\u0063\u0072\u0069\u0070ti\u006f\u006e"
	}
	_eec := _ag.Sprintf("\u004ci\u0063e\u006e\u0073\u0065\u0020\u0049\u0064\u003a\u0020\u0025\u0073\u000a", _badg.LicenseId)
	_eec += _ag.Sprintf("\u0043\u0075s\u0074\u006f\u006de\u0072\u0020\u0049\u0064\u003a\u0020\u0025\u0073\u000a", _badg.CustomerId)
	_eec += _ag.Sprintf("\u0043u\u0073t\u006f\u006d\u0065\u0072\u0020N\u0061\u006de\u003a\u0020\u0025\u0073\u000a", _badg.CustomerName)
	_eec += _ag.Sprintf("\u0054i\u0065\u0072\u003a\u0020\u0025\u0073\n", _badg.Tier)
	_eec += _ag.Sprintf("\u0043r\u0065a\u0074\u0065\u0064\u0020\u0041\u0074\u003a\u0020\u0025\u0073\u000a", _ce.UtcTimeFormat(_badg.CreatedAt))
	if _badg.ExpiresAt.IsZero() {
		_eec += "\u0045x\u0070i\u0072\u0065\u0073\u0020\u0041t\u003a\u0020N\u0065\u0076\u0065\u0072\u000a"
	} else {
		_eec += _ag.Sprintf("\u0045x\u0070i\u0072\u0065\u0073\u0020\u0041\u0074\u003a\u0020\u0025\u0073\u000a", _ce.UtcTimeFormat(_badg.ExpiresAt))
	}
	_eec += _ag.Sprintf("\u0043\u0072\u0065\u0061\u0074\u006f\u0072\u003a\u0020\u0025\u0073\u0020<\u0025\u0073\u003e\u000a", _badg.CreatorName, _badg.CreatorEmail)
	return _eec
}

func (_acb defaultStateHolder) loadState(_bdg string) (reportState, error) {
	_cdf := _aegb()
	if len(_cdf) == 0 {
		return reportState{}, _gb.New("\u0068\u006fm\u0065\u0020\u0064i\u0072\u0020\u006e\u006f\u0074\u0020\u0073\u0065\u0074")
	}
	_bcg := _ed.Join(_cdf, "\u002eu\u006e\u0069\u0064\u006f\u0063")
	_fde := _cd.MkdirAll(_bcg, 0777)
	if _fde != nil {
		return reportState{}, _fde
	}
	if len(_bdg) < 20 {
		return reportState{}, _gb.New("i\u006e\u0076\u0061\u006c\u0069\u0064\u0020\u006b\u0065\u0079")
	}
	_bec := []byte(_bdg)
	_bga := _abb.Sum512_256(_bec[:20])
	_ggb := _agc.EncodeToString(_bga[:])
	_afc := _ed.Join(_bcg, _ggb)
	_fce,
		_fde := _fa.ReadFile(_afc)
	if _fde != nil {
		if _cd.IsNotExist(_fde) {
			return reportState{}, nil
		}
		_eb.Log.Error("\u0045R\u0052\u004f\u0052\u003a\u0020\u0025v", _fde)
		return reportState{}, _gb.New("\u0069\u006e\u0076a\u006c\u0069\u0064\u0020\u0064\u0061\u0074\u0061")
	}
	const _ccfc = "\u0068\u00619\u004e\u004b\u0038]\u0052\u0062\u004c\u002a\u006d\u0034\u004c\u004b\u0057"
	_fce,
		_fde = _eae([]byte(_ccfc), _fce)
	if _fde != nil {
		return reportState{}, _fde
	}
	var _cfb reportState
	_fde = _dg.Unmarshal(_fce, &_cfb)
	if _fde != nil {
		_eb.Log.Error("\u0045\u0052\u0052OR\u003a\u0020\u0049\u006e\u0076\u0061\u006c\u0069\u0064\u0020\u0064\u0061\u0074\u0061\u003a\u0020\u0025\u0076", _fde)
		return reportState{}, _gb.New("\u0069\u006e\u0076a\u006c\u0069\u0064\u0020\u0064\u0061\u0074\u0061")
	}
	return _cfb,
		nil
}

func _bad(_gd string, _dd string, _ded string) (string, error) {
	_ebag := _ab.Index(_ded, _gd)
	if _ebag == -1 {
		return "", _ag.Errorf("\u0068\u0065a\u0064\u0065\u0072 \u006e\u006f\u0074\u0020\u0066\u006f\u0075\u006e\u0064")
	}
	_gdd := _ab.Index(_ded, _dd)
	if _gdd == -1 {
		return "", _ag.Errorf("\u0066\u006fo\u0074\u0065\u0072 \u006e\u006f\u0074\u0020\u0066\u006f\u0075\u006e\u0064")
	}
	_abg := _ebag + len(_gd) + 1
	return _ded[_abg : _gdd-1],
		nil
}

func (_egc *LicenseKey) TypeToString() string {
	if _egc._cfd {
		return "M\u0065t\u0065\u0072\u0065\u0064\u0020\u0073\u0075\u0062s\u0063\u0072\u0069\u0070ti\u006f\u006e"
	}
	if _egc.Tier == LicenseTierUnlicensed {
		return "\u0055\u006e\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u0064"
	}
	if _egc.Tier == LicenseTierCommunity {
		return "\u0041\u0047PL\u0076\u0033\u0020O\u0070\u0065\u006e\u0020Sou\u0072ce\u0020\u0043\u006f\u006d\u006d\u0075\u006eit\u0079\u0020\u004c\u0069\u0063\u0065\u006es\u0065"
	}
	if _egc.Tier == LicenseTierIndividual || _egc.Tier == "\u0069\u006e\u0064i\u0065" {
		return "\u0043\u006f\u006dm\u0065\u0072\u0063\u0069a\u006c\u0020\u004c\u0069\u0063\u0065\u006es\u0065\u0020\u002d\u0020\u0049\u006e\u0064\u0069\u0076\u0069\u0064\u0075\u0061\u006c"
	}
	return "\u0043\u006fm\u006d\u0065\u0072\u0063\u0069\u0061\u006c\u0020\u004c\u0069\u0063\u0065\u006e\u0073\u0065\u0020\u002d\u0020\u0042\u0075\u0073\u0069ne\u0073\u0073"
}

func (_agfa defaultStateHolder) updateState(_ga, _acg, _bc string, _bbg int, _fgd bool, _aefa int, _bba int, _feaa _gg.Time, _bda map[string]int) error {
	_bfga := _aegb()
	if len(_bfga) == 0 {
		return _gb.New("\u0068\u006fm\u0065\u0020\u0064i\u0072\u0020\u006e\u006f\u0074\u0020\u0073\u0065\u0074")
	}
	_cbef := _ed.Join(_bfga, "\u002eu\u006e\u0069\u0064\u006f\u0063")
	_adb := _cd.MkdirAll(_cbef, 0777)
	if _adb != nil {
		return _adb
	}
	if len(_ga) < 20 {
		return _gb.New("i\u006e\u0076\u0061\u006c\u0069\u0064\u0020\u006b\u0065\u0079")
	}
	_fbge := []byte(_ga)
	_fgc := _abb.Sum512_256(_fbge[:20])
	_fef := _agc.EncodeToString(_fgc[:])
	_bfa := _ed.Join(_cbef, _fef)
	var _ceg reportState
	_ceg.Docs = int64(_bbg)
	_ceg.NumErrors = int64(_bba)
	_ceg.LimitDocs = _fgd
	_ceg.RemainingDocs = int64(_aefa)
	_ceg.LastWritten = _gg.Now().UTC()
	_ceg.LastReported = _feaa
	_ceg.Instance = _acg
	_ceg.Next = _bc
	_ceg.Usage = _bda
	_gdc,
		_adb := _dg.Marshal(_ceg)
	if _adb != nil {
		return _adb
	}
	const _abf = "\u0068\u00619\u004e\u004b\u0038]\u0052\u0062\u004c\u002a\u006d\u0034\u004c\u004b\u0057"
	_gdc,
		_adb = _cge([]byte(_abf), _gdc)
	if _adb != nil {
		return _adb
	}
	_adb = _fa.WriteFile(_bfa, _gdc, 0600)
	if _adb != nil {
		return _adb
	}
	return nil
}

func (_da *LicenseKey) getExpiryDateToCompare() _gg.Time {
	if _da.Trial {
		return _gg.Now().UTC()
	}
	return _ce.ReleasedAt
}

var _bgb map[string]int
