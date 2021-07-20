package ini_test

import (
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/ini.v1"
)

type testData struct {
	Value1 string `ini:"value1"`
	Value2 string `ini:"value2"`
	Value3 string `ini:"value3"`
}

func TestMultiline(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping testing on Windows")
	}

	Convey("Parse Python-style multiline values", t, func() {
		path := filepath.Join("testdata", "multiline.ini")
		f, err := ini.LoadSources(ini.LoadOptions{
			AllowPythonMultilineValues: true,
			ReaderBufferSize:           64 * 1024,
		}, path)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(len(f.Sections()), ShouldEqual, 1)

		defaultSection := f.Section("")
		So(f.Section(""), ShouldNotBeNil)

		var testData testData
		err = defaultSection.MapTo(&testData)
		So(err, ShouldBeNil)
		So(testData.Value1, ShouldEqual, "some text here\nsome more text here\n\nthere is an empty line above and below\n")
		So(testData.Value2, ShouldEqual, "there is an empty line above\nthat is not indented so it should not be part\nof the value")
		So(testData.Value3, ShouldEqual, `.

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Eu consequat ac felis donec et odio pellentesque diam volutpat. Mauris commodo quis imperdiet massa tincidunt nunc. Interdum velit euismod in pellentesque. Nisl condimentum id venenatis a condimentum vitae sapien pellentesque. Nascetur ridiculus mus mauris vitae. Posuere urna nec tincidunt praesent semper feugiat. Lorem donec massa sapien faucibus et molestie ac feugiat sed. Ipsum dolor sit amet consectetur adipiscing elit. Enim sed faucibus turpis in eu mi. A diam sollicitudin tempor id. Quam nulla porttitor massa id neque aliquam vestibulum morbi blandit.

Lectus sit amet est placerat in egestas. At risus viverra adipiscing at in tellus integer. Tristique senectus et netus et malesuada fames ac. In hac habitasse platea dictumst. Purus in mollis nunc sed. Pellentesque sit amet porttitor eget dolor morbi. Elit at imperdiet dui accumsan sit amet nulla. Cursus in hac habitasse platea dictumst. Bibendum arcu vitae elementum curabitur. Faucibus ornare suspendisse sed nisi lacus. In vitae turpis massa sed. Libero nunc consequat interdum varius sit amet. Molestie a iaculis at erat pellentesque.

Dui faucibus in ornare quam viverra orci sagittis eu. Purus in mollis nunc sed id semper. Sed arcu non odio euismod lacinia at. Quis commodo odio aenean sed adipiscing diam donec. Quisque id diam vel quam elementum pulvinar. Lorem ipsum dolor sit amet. Purus ut faucibus pulvinar elementum integer enim neque volutpat ac. Fermentum posuere urna nec tincidunt praesent semper feugiat nibh sed. Gravida rutrum quisque non tellus orci. Ipsum dolor sit amet consectetur adipiscing elit pellentesque habitant. Et sollicitudin ac orci phasellus egestas tellus rutrum tellus pellentesque. Eget gravida cum sociis natoque penatibus et magnis. Elementum eu facilisis sed odio morbi quis commodo. Mollis nunc sed id semper risus in hendrerit gravida rutrum. Lorem dolor sed viverra ipsum.

Pellentesque adipiscing commodo elit at imperdiet dui accumsan sit amet. Justo eget magna fermentum iaculis eu non diam. Condimentum mattis pellentesque id nibh tortor id aliquet lectus. Tellus molestie nunc non blandit massa enim. Mauris ultrices eros in cursus turpis. Purus viverra accumsan in nisl nisi scelerisque. Quis lectus nulla at volutpat. Purus ut faucibus pulvinar elementum integer enim. In pellentesque massa placerat duis ultricies lacus sed turpis. Elit sed vulputate mi sit amet mauris commodo. Tellus elementum sagittis vitae et. Duis tristique sollicitudin nibh sit amet commodo nulla facilisi nullam. Lectus vestibulum mattis ullamcorper velit sed ullamcorper morbi tincidunt ornare. Libero id faucibus nisl tincidunt eget nullam. Mattis aliquam faucibus purus in massa tempor. Fames ac turpis egestas sed tempus urna. Gravida in fermentum et sollicitudin ac orci phasellus egestas.

Blandit turpis cursus in hac habitasse. Sed id semper risus in. Amet porttitor eget dolor morbi non arcu. Rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Ut morbi tincidunt augue interdum velit. Lorem mollis aliquam ut porttitor leo a. Nunc eget lorem dolor sed viverra. Scelerisque mauris pellentesque pulvinar pellentesque. Elit at imperdiet dui accumsan sit amet. Eget magna fermentum iaculis eu non diam phasellus vestibulum lorem. Laoreet non curabitur gravida arcu ac tortor dignissim. Tortor pretium viverra suspendisse potenti nullam ac tortor vitae purus. Lacus sed viverra tellus in hac habitasse platea dictumst vestibulum. Viverra adipiscing at in tellus. Duis at tellus at urna condimentum. Eget gravida cum sociis natoque penatibus et magnis dis parturient. Pharetra massa massa ultricies mi quis hendrerit.

Mauris pellentesque pulvinar pellentesque habitant morbi tristique. Maecenas volutpat blandit aliquam etiam. Sed turpis tincidunt id aliquet. Eget duis at tellus at urna condimentum. Pellentesque habitant morbi tristique senectus et. Amet aliquam id diam maecenas. Volutpat est velit egestas dui id. Vulputate eu scelerisque felis imperdiet proin fermentum leo vel orci. Massa sed elementum tempus egestas sed sed risus pretium. Quam quisque id diam vel quam elementum pulvinar etiam non. Sapien faucibus et molestie ac. Ipsum dolor sit amet consectetur adipiscing. Viverra orci sagittis eu volutpat. Leo urna molestie at elementum. Commodo viverra maecenas accumsan lacus. Non sodales neque sodales ut etiam sit amet. Habitant morbi tristique senectus et netus et malesuada fames. Habitant morbi tristique senectus et netus et malesuada. Blandit aliquam etiam erat velit scelerisque in. Varius duis at consectetur lorem donec massa sapien faucibus et.

Augue mauris augue neque gravida in. Odio ut sem nulla pharetra diam sit amet nisl suscipit. Nulla aliquet enim tortor at auctor urna nunc id. Morbi tristique senectus et netus et malesuada fames ac. Quam id leo in vitae turpis massa sed elementum tempus. Ipsum faucibus vitae aliquet nec ullamcorper sit amet risus nullam. Maecenas volutpat blandit aliquam etiam erat velit scelerisque in. Sagittis nisl rhoncus mattis rhoncus urna neque viverra justo. Massa tempor nec feugiat nisl pretium. Vulputate sapien nec sagittis aliquam malesuada bibendum arcu vitae elementum. Enim lobortis scelerisque fermentum dui faucibus in ornare. Faucibus ornare suspendisse sed nisi lacus. Morbi tristique senectus et netus et malesuada fames. Malesuada pellentesque elit eget gravida cum sociis natoque penatibus et. Dictum non consectetur a erat nam at. Leo urna molestie at elementum eu facilisis sed odio morbi. Quam id leo in vitae turpis massa. Neque egestas congue quisque egestas diam in arcu. Varius morbi enim nunc faucibus a pellentesque sit. Aliquet enim tortor at auctor urna.

Elit scelerisque mauris pellentesque pulvinar pellentesque habitant morbi tristique. Luctus accumsan tortor posuere ac. Eu ultrices vitae auctor eu augue ut lectus arcu bibendum. Pretium nibh ipsum consequat nisl vel pretium lectus. Aliquam etiam erat velit scelerisque in dictum. Sem et tortor consequat id porta nibh venenatis cras sed. A scelerisque purus semper eget duis at tellus at urna. At auctor urna nunc id. Ornare quam viverra orci sagittis eu volutpat odio. Nisl purus in mollis nunc sed id semper. Ornare suspendisse sed nisi lacus sed. Consectetur lorem donec massa sapien faucibus et. Ipsum dolor sit amet consectetur adipiscing elit ut. Porta nibh venenatis cras sed. Dignissim diam quis enim lobortis scelerisque. Quam nulla porttitor massa id. Tellus molestie nunc non blandit massa.

Malesuada fames ac turpis egestas. Suscipit tellus mauris a diam maecenas. Turpis in eu mi bibendum neque egestas. Venenatis tellus in metus vulputate eu scelerisque felis imperdiet. Quis imperdiet massa tincidunt nunc pulvinar sapien et. Urna duis convallis convallis tellus id. Velit egestas dui id ornare arcu odio. Consectetur purus ut faucibus pulvinar elementum integer enim neque. Aenean sed adipiscing diam donec adipiscing tristique. Tortor aliquam nulla facilisi cras fermentum odio eu. Diam in arcu cursus euismod quis viverra nibh cras.

Id ornare arcu odio ut sem. Arcu dictum varius duis at consectetur lorem donec massa sapien. Proin libero nunc consequat interdum varius sit. Ut eu sem integer vitae justo. Vitae elementum curabitur vitae nunc. Diam quam nulla porttitor massa. Lectus mauris ultrices eros in cursus turpis massa tincidunt dui. Natoque penatibus et magnis dis parturient montes. Pellentesque habitant morbi tristique senectus et netus et malesuada fames. Libero nunc consequat interdum varius sit. Rhoncus dolor purus non enim praesent. Pellentesque sit amet porttitor eget. Nibh tortor id aliquet lectus proin nibh. Fermentum iaculis eu non diam phasellus vestibulum lorem sed.

Eu feugiat pretium nibh ipsum consequat nisl vel pretium lectus. Habitant morbi tristique senectus et netus et malesuada fames ac. Urna condimentum mattis pellentesque id. Lorem sed risus ultricies tristique nulla aliquet enim tortor at. Ipsum dolor sit amet consectetur adipiscing elit. Convallis a cras semper auctor neque vitae tempus quam. A diam sollicitudin tempor id eu nisl nunc mi ipsum. Maecenas sed enim ut sem viverra aliquet eget. Massa enim nec dui nunc mattis enim. Nam aliquam sem et tortor consequat. Adipiscing commodo elit at imperdiet dui accumsan sit amet nulla. Nullam eget felis eget nunc lobortis. Mauris a diam maecenas sed enim ut sem viverra. Ornare massa eget egestas purus. In hac habitasse platea dictumst. Ut tortor pretium viverra suspendisse potenti nullam ac tortor. Nisl nunc mi ipsum faucibus. At varius vel pharetra vel. Mauris ultrices eros in cursus turpis massa tincidunt.`)
	})
}
