package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _esc_localFS struct{}

var _esc_local _esc_localFS

type _esc_staticFS struct{}

var _esc_static _esc_staticFS

type _esc_file struct {
	compressed string
	size       int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_esc_localFS) Open(name string) (http.File, error) {
	f, present := _esc_data[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_esc_staticFS) Open(name string) (http.File, error) {
	f, present := _esc_data[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		gr, err = gzip.NewReader(bytes.NewBufferString(f.compressed))
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (f *_esc_file) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_esc_file
	}
	return &httpFile{
		Reader:    bytes.NewReader(f.data),
		_esc_file: f,
	}, nil
}

func (f *_esc_file) Close() error {
	return nil
}

func (f *_esc_file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_esc_file) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_esc_file) Name() string {
	return f.name
}

func (f *_esc_file) Size() int64 {
	return f.size
}

func (f *_esc_file) Mode() os.FileMode {
	return 0
}

func (f *_esc_file) ModTime() time.Time {
	return time.Time{}
}

func (f *_esc_file) IsDir() bool {
	return f.isDir
}

func (f *_esc_file) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _esc_local
	}
	return _esc_static
}

var _esc_data = map[string]*_esc_file{

	"/static/main.css": {
		local: "resources/static/main.css",
		size:  224,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xffl\xce\xc1\x0e\x820\f\x06\xe0\xfb\x9e\xa2r\xc7輍W\xf1\xd2\xc1\x18\x8b\xd0\xceZ\x8d\xc4\xf0\xee\x12\x90\x84\x83\u05ef\xfd\xff\xd6s3\xc2\xc7\x00dl\x9aD\xb1T\xce\x0e\xec)\xbf\xab\x1dzV\xe5a\xf3ɘ#\xe1ˣ,\xc9\x01%&\xfa\xb3\x13\xfb1w\xa9f*\x89)8\x1fZ\x96\xb0D" +
			"f\xd3@꠸ڳ\xb5E\xb5X\xcf\xe2@\x05\xe9\x91Q\xe61\x1cҐY\x14I\xd7\xc2\xfb3շ\xfd\xcd>\xb4s\xc9e}\xf6g\x92b\xb7\xe1d\xbe\x01\x00\x00\xff\xff 9Z\xfa\xe0\x00\x00\x00",
	},

	"/static/main.js": {
		local: "resources/static/main.js",
		size:  459,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xfftP\xcbj\xc30\x10<\xdb_\xb1(\xa1^Ac\x97BO\xb1si?\xa1\xb7҃,\xafbQE\n\x96\x926\x94\xfc{%\xb9\x04\xd2\a\x06k\xb5;\xb33\xa3%\x0eN\x1evd\x03\xaf'\x12\xc3\t\xd5\xc1ʠ\x9dE\x0e\x9feY4\r\xac~~\xb9\xfb,zC\xe0\x14HgC\xe4\xfb" +
			"\u007f\xb0e\xa1\x15\xe0\x12\x81-z7\x9c\xc0\x8a#\x03^\x8fag\x90\xd7a\xd2;\xe4I\xaa(~azRn\"\x84\xaa\xf5{a7\x8f\xdfBm\x93\xaf7\xb6\xf7\xfbuu\v\xd5K+@\x0f\x1d\xf3\xa3{_\x8dz\xa0Up\x92\xc18\x91\xea\u0602mR\xabm\xc4\xe6\xb5\x02\xbe.\x8bs\xf4\xb4D\xb6\xb8\x86\xf3Z\x1a-\xdf." +
			"\xf9\x01阞%[\xcbe\xbd\x9f\xf2\xf9DJ\x1cL\xc0\xb4+/\xbax\x8ey\xdcvk\b\xef\x1f\xee\xf2t\x8e\x1eF\xed\xe3\x88>\"\a\xba\xae\x03\x96Tټ\xba\xb8\x9a\xe7\f,s\xcf@\xc6\xd3_\x98\x99\x9d11M,\xca\xf4\xfb\n\x00\x00\xff\xff>\xe9\x83`\xcb\x01\x00\x00",
	},

	"/templates/_base.html": {
		local:      "resources/templates/_base.html",
		size:       91,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xb2Qt\xf1w\x0e\x89\fpU\xc8(\xc9ͱ\xe3\xb2\x01Q\n9\x89y\xe9\xb6J\xa9yJv\\\\\n\n\xd5\xd5%\xa9\xb9\x059\x89%\xa9\nJ\x19\xa9\x89)J\nz\xb5\xb5h\xe2I\xf9)\x95\x10q.\x1b}\x88Q\x80\x00\x00\x00\xff\xffmS\aw[\x00\x00\x00",
	},

	"/templates/_body.html": {
		local: "resources/templates/_body.html",
		size:  1587,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xbcU\xc1\x8e\x9c0\f=\x97\xaf\x88\xd23D۽U\x80T\xa9\x1f\xd0C\xef\x95!\x19\xc8*$lb\xaa\x19\xa1\xfd\xf7:\x81a\aڙ\xbdu\xa5\x15\x8ec?;\xcf/\x99y\x96꤭b\xbcq\xf2\xc2\xdf\u07b2\xac\x8cV\x9d1VJ\xfd\x9b\xb5\x06B\xa8x\xeb,\x02\xc5y^g\xb4E" +
			"\x9b\x16\xb6M2\x1b\xf0l\xf9\xe4\x04\b\x93A^\xa7\xb8;0\xf9\xc9LZn1\xfb\xa8\x15\xa8W SA\xb6\xfd\x95̈́\xe8,\xc3˨*\xbe,\xf8!\r]\xd7\x19\xc5Zg\f\x8cAI\xce$ \xac\xee\xd8\xc2⿺\xc1w\n+\xfey\xc9\xe6\f\xbc\x86\\\x9dG\xb0RɊ\x9f\xc0\xc4\xd8\xe4\x8d\xdd{g\xb6R\xbb\xd6" +
			"\xa8\xb9@I\xd7f\x82ϝ5\x17^\xff\\ڡ\f\xdd\x01jgK\x11\xe3\x1e\xa4j\xaa\x93'\xf8\xff\x15Z\x8a\x85ʝ\x0f\x0e\xbc6\x9e(\xe1\xac\xf7\xeaTqA0\xb7蝹\x8c},\xc16+\xef\xddpe\xae\xd7R*[q\xf4\x93\xda\x1a(\x05\xdc\xcc_\x90\x00\x0er\xd0rc\xfa\xd0\xcbu\x88۔\xf7*\x99\xccM" +
			"\xfcU\x97\xf49\x0e\xcc蚎\xb9\x9e\xe8\x97\x18\xa1S\x81\xd7ߌa?\xa2\x19\x1b,\x05\x05=̒\xca($\x99\xd5\xdf\x17\xe3~n)&Sg\x9f\xc889?\x1c\x8e\x94\\\xab\xedu\xd7#g\xa45\x12lP\xe0۞\x88l\xa3zR\xc9\xd5u\xe8\xeb\xe6\x06E\xb0\xbc\xf3n\x1ay*\x98\xf6\xb5\x1d'\\\xaf\x0e\xaa3\x15" +
			"\xb00\x90\xfd\xcawi\xabʏ\xe8\xfb\xf9$\xd7\xee2\x86\xa9\x194nP\rZF\xff\xef\xaf\xc1\x87r\xd9\xce\xf9@0\xff\x90\xa9\x88M\xdf\xd1\xd1͢\x14\xc4\xed\xfav\xcd3\xaaa4\x80\xf4\xecř\xe7\v\xb5\x81\xb3\x82\x1e\xc0cD\xe4CY\\6\xb3\r4Z\xa1\xf5zD\x16|[\xf1\x1eq\f_\x85\x80\x178\x17\x9ds" +
			"t\xe9aԡhݐ|\xa4\x85&\x88\x97\xd7I\xf9\x8bx*\x9e\x9e\x8a/\xeb\xaa\x18\xb4-^B:f\x02\xac\xefa\x0fpn\xa5-\x1a\xe70\xa0\x871.\"\xfe\xe6\x10\xcf\xc5s\xc4\rﮏѣ\xa0\x90\x9e\xa6\x96\xf0\xff\x8a͈\xf4\xf4{\x90ͳ\xb2\x92(\xf8\x13\x00\x00\xff\xffX\xcb&\x1d3\x06\x00\x00",
	},

	"/templates/_delete.html": {
		local: "resources/templates/_delete.html",
		size:  880,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\x84Sˮ\x9b0\x10]\xdf|\xc5ȋ\xee(\xaatw5d\x91,[\xb5R\xf3\x03\x06\x0f`\xc9x\"3\xb4E\x88\u007f\xafy5$A7\v\x84\x99ǙsΘ\xbe\xd7X\x18\x87 4Zd\x8cj\xd2ʊa8\x1c\xa46\xbf!\xb7\xaai\x121E\xa1P\x1a\x05\x18\x9d\x88\x9c\\a|" +
			"\x1d\xcdM\x02Xe\xc6i\xfc\x9b\x88\xe8\x8b\x00O\x16\x13\xa1\x8d\xb2T\nPިȪ\f\xadE\x9du\x01\xac\xfb>\xc2}\x1bCK\xba2Z\xa3K\x04\xfb\x16Ez\x00x\x1a\x1e-pcr/\x1d\x181:\x0e\xf9\xa9`\xaf\xa4\xc2\xc0\xdf/\bSI\xd62\x93\x03\ueb81\xef\xfc!֞\xdcR\x13\x94i\xc5*\xccnj\xf3" +
			"\x1fh\x97\xf3'656_e<\xc3l\x86T\xef\xf74ذ]l\xbcs\"=ͦ\xc2y2U\xc6\xd5\xfb\x8a\"\xe3 \xe6\x03e\x19\xe9n\xab뚞\t:j\xe1\x8fr\fL0\xef\tdÞ\\\x99\xf6\xfd\xe7\xcb\xc8b\x18d\xbc\x84\x8e2\xben\x11\xd6\x01\x1a\xb3\xb6\x8cZ\x1f\xf8m*^\x11*\x88x\xb2\xfaM\x16" +
			"\x14$\xa9\x9c\r\x05\xaf\xe2\xdb\xe4\xe3\x12[\xafP\x8d\\Q\xf0\xe4\xe7\x8f_\x97\x8d\x96\x17[\xca\xd8Ax\xc2E,Tky\u007f_\xe9I\xb9\x1c\xed\xf3nn\xe0\x8fpʕ\xe8a~\x89t\xdd\xc8\xda\xff&\xe3Qփ\x1bw\xc7\xe5\xb0\x1a\xd5\xf7\xe8t\xf8\xab\xfe\x05\x00\x00\xff\xff|\x16\x01\xe3p\x03\x00\x00",
	},

	"/templates/_edit.html": {
		local: "resources/templates/_edit.html",
		size:  1474,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xb4TKo\xdb0\f>ǿB\xd0ή\xb1\xbb\xe3\xc3\xda\x1d\x87\x05H\xb1\xbbbю0=<\x99\xce\x10\x18\xfe\xef\xa3\x1eNR,\x1d\xb6\x0e=Ĉ\xf8\xf8\xf8}$\xc1y\x96\xd0)\v\x8c\x0f\xa2\x87R\xb4\xa8\x9c\x1d\xf9\xb2\x14E-\xd8\xd1C\xb7\xe5\x1f8\x93\x02E\x89\xae\xef5l\xb9q" +
			"R\xe8\xd5&|\x0fH1\xad\xb3\x9d\U000a6520\x01\x81\xb3V\x8bq\xdc\xf2\x03ZF?2wb\xd2ȆI\xebҫ\xfe\x88\xecǤ\xda\xef\xd19\x1aμ\vЇ\t\xd1Y\xde<E\x94\xba\x12͕F5\xcf\x0f\xcf\n5,\xcb\xff\xc2\u007fS\xf0\x93\xedHp\xacP\xcc3X\x19$\x17\xd7v\x90\x1e\x04\x8b\xa9\x13Ǐ\xcd" +
			"\xb58\xabG#\xb4\x0e\x96\xfdt\xc8ƺJƺ\xa2X\xca\xe8\x9c7\xcc\x00\x1e\x9d\xdc\xf2\xdd\xd7\xfd3o\n\xc6\xe6Yu\xec\xe1\xb3T\xb8,TT\x8f\x94I\xe6Z\xd9aB\x86\xe7\x818\x1e\x95\x94`9\xb3\xc2\x04\xc6N\x9e9;\t=у*~\xa2\xf7\xdeM\xbe\r]h\xfe\x9ck`\x1cI\xe3m\xfa\xa33F\xe1\x97\xe4" +
			"\xc8\b\x17\xf5\x04&\xd5im\xed ,h\x16\xbfk{c\xf8\x9d\xa82\xb2l\x8a\xe8%\u007fjs\xe64N\a*\xf9\xdb\xc4\x06\xaf\x8c\xf0\xe7\x95kZ\xbc\v\xd5Q\x9c\x807{\xfa\xd6UBkV\xf0\u007fX\a\xde<\nۂNS\xce\xf97\xe4Ch\xef\xdd4\x10\xf7\xcd\xedl\x8a\xcdK\ryo^ۺ@>\xf0\r\xc9W" +
			"\xbe\x9b˄7\u007fՑ\x95\xf4\xfd\x8e\x00a߫\x90F\x97\xd9\xef<\x9ch\xb5\xdf, 翗\x86!\xc1\xbfR煒'\xd5uo\x96\x11\x92\xdfK\x83$\xec{\x15\"\xfb\xbcb\x15\xedX^\xb8\xfc\xff\xd68\xcf\bf\xd0\x02\xe9ʄ\xa9\x96\xeb\xa9a\x0f\xf1\xd8T\xe1vīt\x8dKW\xb5̗7ƭ5\u007f\x05" +
			"\x00\x00\xff\xffp\xe4\xf6e\xc2\x05\x00\x00",
	},

	"/templates/_head.html": {
		local: "resources/templates/_head.html",
		size:  413,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff|\x90OO3!\x10\xc6\xef\xfb)\b\xe7\x97%\xaf\xbd\x18\xb341ƃg5\xf1f(Le\"\v\xeb2m5\x84\xef.t\x9b\xb6'o3\xcf3\xbf\xf9\x97\xb3\x85-\x06`܁\xb6\xbc\x94\xae\x1bZ\xb4\xee\x18\x1bF ͌\xd3s\x02R|G[q\xcb/\x86#\x9a\x04|\xedp\xaf\xf8" +
			"\x9bx\xbd\x17\x0fq\x9c4\xe1\xc6\x03g&\x06\x82P\xa9\xa7G\x05\xf6\x03\xae\xb8\xa0GP|\x8fp\x98\xe2LW\xa5\a\xb4䔅=\x1a\x10\xc7\xe4\x1fÀ\x84ڋd\xb4\a\xf5\xbf\xb6i}\b\xc9\xc3:\xe7\xfe\xa5\x05\xa5\xe4|@r\xac\u007f\xdemN\n\x13\xac\xdá`K\x19\xe4\x82\x1ci\x8f\xe1\x93\xcd\xe0\x15O\xf4\xe3!" +
			"9\x80\xba\x86\x9ba\xabx;*\xddI9\xeaocC\xbf\x89\x91\x12\xcdzj\x89\x89\xa3<\vrկ\xfa\x1biR\xbah\xfd\x88\xb5*\xa5\xe5ؿ\xc6\xc8w\x99\xa8\xfe\xca\xd4Agf\x90\xcb\xe7\xbb\xd3\xd2\xddo\x00\x00\x00\xff\xff<\x87=j\x9d\x01\x00\x00",
	},

	"/templates/all-pages.html": {
		local: "resources/templates/all-pages.html",
		size:  285,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xffL\x8f=\x8e\xc4 \f\x85{Na\xd1'hەC\xb77\xd8\v\xa0\xe0$H\bF\x86\x99\x06q\xf7\x01\x92\xf9i0O\xef\xe9{v)\x966\x17\b\xe4\xcd\xec4\x995\xbb\x18\x92\xacU\x94B\xc1\xb6)>\x915\x86L!wW\xe0\xf1\xa3K\x99\xff]\xf6T+\xaa&{\xd4m0\xff1G" +
			"\x1e\x19\xeb\x1e\xb0z\x93\xd2\"\x8d'\xce0\xdeɚ\xb0\x13K\xe0\xe8\xe9r\xa4\x16\x00\x982ǰ\xeb\x01@u\xa9_h=/&\xaa\xc6\x1cM\xe4\x13\x8d\x92\xbb\xd7Mrg¹O\xeaA\xef4\x1a8\x98\xb6E\xaaF\xa8U\xea1P\x19\x8d\xaa\xd9\xef\x13Qu\xc6\xf7\xc5\xe7\xe7\x19\x00\x00\xff\xff\x11\xc7\xf8P\x1d\x01\x00\x00",
	},

	"/templates/deleted.html": {
		local: "resources/templates/deleted.html",
		size:  307,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xffD\x90An\xc4 \fE\xf7\x9c\xc2b?\x83\xba\xad\x1c\xba\xea\rz\x01\x948\t\x12\x82\xca\xd0J\x95\xc5\xddkȌf\x03\x98\xffy\x1f[d\xa3=f\x02\xfb\x1d\x0e\xba\x85\xb5Œ\xab\xed݈P\xdet7\xe6\xe5YKn\x94ې\r\x9eo^\xe4\xfe\x15[\xa2\xde\xd1i9\xacq\x87\xfb'" +
			"s\xe1\xe9\xd9\xe2/\xac)Ժؐ\x88\x1b\xcc\xf5\xb6\x85|\x10[\xe0\x92\xe8\xa1Xo\x00\xb06.\xf9\xf0\x13\x80\xeeQ\xbd\x83\xe6<\x99\xe8\x949\x93(U\x9a!?i\xbc\x15\xe1A\x85\xebGU\x15\xe5\xa5\xe81\xc0ɴ/\xd6)\xa5\xf7\x8f\xab\xc5匵\x15\xfe\xb3~ޢ\v\x1e\x9d\xba'\xe8j\x1c\xdd\x00\xbf\xe6\xf0<\xfc" +
			"\a\x00\x00\xff\xff\x87\xaf\x8b\xb23\x01\x00\x00",
	},

	"/templates/diff.html": {
		local:      "resources/templates/diff.html",
		size:       82,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xaa\xaeVHIM\xcb\xccKUPJM\xc9,\xd1M\xce\xcf+I\xcd+QR\xa8\xad\xe5\xe2\xb2I\xc9,S\xc8L\xb1UJ\xcaO\xa9T\xb2\xb3)(J\xb5\xab\xae\xd6s\x02\xf2jkm\xf4A\\\x1b}\xa0\x1a;..\xa09\xa9y) ]\x80\x00\x00\x00\xff\xff\x13!?nR\x00\x00\x00",
	},

	"/templates/edit.html": {
		local: "resources/templates/edit.html",
		size:  243,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xffl\x8f1n\xc4 \x10E{N1\x9a>\xb1R\xa4\xc3.\x92:UN@`\x9c \x01c\xc1؉\x85\xb8{X{w\xab\xed\x90\x1e\xff=M\xad\x8ef\x9f\b\x90\x9c\x97'\xcbI(\t\xb6\xa6\x94\x16\xfa\x13\x93\xc9@2\x91F\xfcb\xb7#d\xfe-#\xbe\xbc\"\xd8`J\u007fΜ\xe3\xb1\xcb\x1c" +
			"p\xaa\xf5\xf9\xad\xff\xfb\xe45[jM\x0f7ɤ\xb4\xf3ۤ\x00\xb4O\xcb* \xfbҥ\x17\x8c\xd7@\xa4R\xcc7=6\xc3\x12\x8c\xa5\x1f\x0e\x8e\xf2\x88\xef\x1c\xa3\x17\xb8/6\x13\xd6n\xe8\xf5\x93|\x9c\xa05\xec\xdd\xe1\b\xabZ)\xb9~\xd8\u007f\x00\x00\x00\xff\xff뾉)\xf3\x00\x00\x00",
	},

	"/templates/history.html": {
		local: "resources/templates/history.html",
		size:  719,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\x8c\x921\xef\xd30\x10\xc5\xf7|\n\xcbB\b\x864bE\xa9\x19(\x03\x03\x12\x03bw\xe3kb\xe1\xd8\u007f\xecKQe\xe5\xbbs>\xe3\xb6H\f\f\xad\xce/\xf6ｳ/g\x03\x17\xebA\xc8\x17=C\xaf'\xb4\xc1'\xb9\xef]7j\xb1D\xb8\x1c\xe5\x90\xf3\xe1\x9bE\a\xfb.\xc5\xe4tJ" +
			"GyF/\xe8\xd7\xd3i\xbd9\x14/\x9bs}\xb4\xf3\x82\xe2\xe7f\xa7\x1f\xfc1\xadR\xc4\xe0\x80\xb6o\x88\xc1K\xf5\xdd\xc2/\xf1\x95\x9c\xc6A\xab\xae\xcb\x19\xbc)^\x8f\x18S\xf0\b\x1ek\x82\xe5\x9dzx\x8b1\xad\xda9\xb5\u0604!\xdeơ.ǁv\x15\x82\xbd\x88ç\x18C\xe4\xa3\xc6^[V\xed \xa2\xe0\xff\xdeh" +
			"?Cl\xb1X\x93\xaa\x13\xc4\xc6\x18\xfc\xac\x18@\xe8\xbaz/Ⱦ1ǁ\x985\xb4K\xc0&\xa8\xcf\x0e\x9a\r/*\f\x17ЦT\xa5\x8e\xb5`Y\x9d4R\xefT<i_ %\xbe\x92&SŇ\x8a\xc2 f\x9e\x83\xb9\x155\xe7X\x9a\x10\x87\x8fa]-&JR\xc4\xd2\xfe\t\x1c \xb0p\xf7\x1dєK,\xc6\xfb" +
			"ND\xf3\x90\xdf\x18\xde\xff\xb6\xa9\xcd\xf7\xde\xe1\xffa\x9e\a\xe5U{\xad\x0fu\x94\x8eWz\xf1\xd7\x11\xae6\x95\x15\x11>\x9fh\x8c\n\xeaOۅ\xa6\xd5?\"\xf0d\xb0R;\xbf\x0f\v)\xe5\xaa\xff\x9e\x9fZ\xfc\x0e\x00\x00\xff\xff\x81\x9c\xe7\x17\xcf\x02\x00\x00",
	},

	"/templates/preview.html": {
		local:      "resources/templates/preview.html",
		size:       67,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xaa\xaeNIM\xcb\xccKUPJM\xc9,\xd1M\xce\xcf+I\xcd+Q\xaa\xad\xe5\xe2\xb2I\xc9,S\xc8L\xb1UJ\xcaO\xa9T\xb2\xab\xae\xd6s\x022jkm\xf4\x81\x12v\\\\\xd5թy)@\x85\x80\x00\x00\x00\xff\xff\x97\xb37~C\x00\x00\x00",
	},

	"/templates/printable.html": {
		local: "resources/templates/printable.html",
		size:  436,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xffTQMO\x021\x10\xbd\xf3+j\xcf\xee6\xc8Ř-\x89\"\aOz\xc0D\x8f\xa5\x1d\xe8\xc4~\xac\xdba\x91\x10\xfe\xbb-\v\n\xa7}\xf3\xe6\xbd7/\xdb\xe6\xe6\xf9u\xb6\xf8|\x9b3K\xdeMGM\xf90\xa7\xc2Zr\b|:b\xac\xb1\xa0L\x01\x19z ŴU]\x02\x92|C\xab" +
			"\xea\x9e_\xae,Q[\xc1\xf7\x06{\xc9?\xaa\xf7\xc7j\x16}\xab\b\x97\x0e8\xd31\x10\x84\xec{\x99K0k\xb8r\x06\xe5A\xf2\x1ea\xdbƎ.\xc4[4d\xa5\x81\x1e5T\xc7\xe1\x96a@B媤\x95\x039\xceAC\x12!9\x98\xee\xf7\xf5\xa2\x80á\x11\x03sZ;\f_\xac\x03'y\xa2\x9d\x83d\x01\xf2%\xdb" +
			"\xc1J\xf2\xd2<=\b\xe1Տ6\xa1^\xc6H\x89:ՖAG/\xfe\b1\xa9'\xf5\x9d\xd0)\xfds\xb5ǬJi(҈\xf3\x1fk\x96\xd1\xec\xce\xd7\xed\xf8\xaaY\x1e\a\xde`\xcf\xd0H^\xb4\xbcH\x9e2(\x8a\xbc8\xe5\r1\xd9s|\xa3\xdf\x00\x00\x00\xff\xff\x82\x16\xf8̴\x01\x00\x00",
	},

	"/templates/search.html": {
		local: "resources/templates/search.html",
		size:  482,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xffd\x91\xcfj\xc4 \x10\x87\xefy\x8a\xc1\xfbFz-nn\xbd\x95B\xff\xbc\x80\xc4I\"\x88\x86\xd1\x14\x8a\xe4\xdd\xeb\xe8f7\xb0\x97Df~\xf3}f\x92\xb3\xc1\xc9z\x04\xb1\xea\x19/zL6\xf8(\xf6\xbd\xcb\x19\xbd)ﮜ\xe0\b\x8d\xc1'\xf4I\x007\xd4\xf22|\xa3\xa6q\x01¸\xb9" +
			"\x14a\n\x04*&\n~\x1er\xee?7\xa4\xbf}W\xf2VR\xb2\x8c0\xd0Nп\x11\x05\xaa\x1cc\u007fat:ƫ\xd0\x0e)A}^\x8c\xf63\x92\x00\n\x0eo\x1d1tp\x17T\xc0\x9d\xfd\nEx0\x95,\xccjB\x17\x91%\xcdٮ\xfb\xd5n[ݛcd\xceĲ\xe7\x80r\xb6*5,\x84\xd3U\xc8\"\xf9" +
			"\xb1\xc9\x15\xa6\x18\x1eg%u\x8d5ډ\xf7^\xd6\xc6\x1c.2\xab\x8cp\xbaQ9ז\\\x9a\x92g[稶\xda\xe9#\xd4:|\x84Ӳ7oz%זi\u007f\xeb\xe9\xf0\x1f\x00\x00\xff\xff\xee\xe0H\x80\xe2\x01\x00\x00",
	},

	"/templates/view.html": {
		local: "resources/templates/view.html",
		size:  1782,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xbcUMo\xdb0\f=/\xbfBP\x81\xdd\\\xa3w'\x03\xf6\x01\xec2\xach\x8b\xddi\x8b\xb1\x85ɖ'\xcbI\x03#\xff}ԗ\xedf\xcbz\x18\xb6C\vY\xd2{|\xa4\xf8\x98i\x12\xb8\x97\x1d2\xdeC\x8d\x19TV\xean\xe0\xe7\xf3fS\by`\x95\x82a\xd8rat/\xf4\xb1c" +
			"\xfd\xa8Tfd\xddX\xf6c\x94\xd5w\xbe\xdb0V\x94\xa3\xb5\xbaK\x97K\xdb1\xfaˈ\x19FeY\x02gV\u05f5B\xe6\x0f\x87\x963{ꑮ{0gR,q\xbe`7\xdeq&\xc0BD-g\x9c\x81\x91\x90\xe1s\x0f\x9d@\x02Y3\xa2\xd7AJ\x06\xdaL:ju\xea\x1bY\x91\xb0y\x95U\xba\x8e\xf8F\n\x81" +
			"]B\x17\xb9C\xfe\x86\xa4\x02\x83v}^\xe4A\xb0_\x8f\xea\xb2BYK\xd293\xdaI\x0ek\x1fNA\x89J\xa1(O\x97YƠJFPop\xc0\u0382{\x89xHǰ\xa2\x94\x16]\xf5\xa0\x94T\x80\xe7-ϨT\x8d\xc1\xfd\x96\xe7\xd3t\xfb$\xad\xc2\xf3\xf9]x\xcb\xedA\xe2\xf1\xed^\x9b\x16\xec\xb67\x92\x88K" +
			"\x853\xef\x9b\xfb\xb4\xc5\x0eh\x06\x02\xa4\x809Da\xb9\x92\xffQ\xa2\x81\xe3\"\xee\x1b\x1d\xb0G=\x9a\n\xff\x9d\xac\x17M\xd6j\x01*큩\xd1n\xf9\r5\xce^\x9a\x96\x1aZ\xa1Ŕ\xc8͢\xf3\xa3?\xb8&\xb1\xc8G\xb5\xdb\x149\xf9iG\xb6\x82\xeb\x85h\xe4`\xb59\xf1kN\xbat\xdf줐at\xd2\xees\xa0" +
			"\xf1B\xfe\x18\x0f\x85\xb4\u007f\x1b\xec\x13q\x84Hӄ\x9dp\x93\x83V,\x8d\x15\xaa\x9d\xa57\xe1̏\x94\xe6n\xb7\xa8\x98&\xb9g\xb7\x0fx\x90\xae\xef\xceg2^\vJ\xb9\x1b\xcb&\x19/nz\xf2\"'\x8a\xcd%rS\xb8\xeea1\xab뙶h\x1bM\x13\xe3\xfe\xeb\xe3S\x98\\\xb2\xebG\x1b\aQ\x98\b\x9cu\xd0җ" +
			"u\x04\x9c\x1d@\x8d\xf4\xb50\xbe\x86+\xb58\xada\x0fp|O[\xaf\x03[\x1c\x06\x9a\xc13\x96\xd2CCW5{Q\x10\xa2q<\xab\xe1L\x83\t\x15\xf3\xff\xd3㥡ry+\xf3\xeafkĹ\x1d\xd4\fc\xd9^o\x87Y\x96\xf1\xb2\xf8.\xc8[O\xc3\xe8\xb6_\xba\xed*\xe5\xee\x03t\x15\xaa\x95a\xbcE\xe6\x05y\xc6" +
			"\xbd몵|B\xee\x87\"\xe4A1Bm\x13`\x9a\xc8\xe0\xbd\x02K\x9d\x17ܚEGߺ\xfeK4?\x03\x00\x00\xff\xff6\xc5MP\xf6\x06\x00\x00",
	},

	"/": {
		isDir: true,
		local: "resources",
	},

	"/static": {
		isDir: true,
		local: "resources/static",
	},

	"/templates": {
		isDir: true,
		local: "resources/templates",
	},
}
