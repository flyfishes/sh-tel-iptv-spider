package model

import "encoding/xml"

const PrefixHeader = `<?xml version="1.0" encoding="UTF-8"?>`

/*
<?xml version="1.0" encoding="UTF-8"?>
<tv generator-info-name="" generator-info-url="">
	<channel id="28">
		<display-name lang="zh">浙江卫视</display-name>
	</channel>
	<programme start="20230223000000 +0800" stop="20230222130000 +0800" channel="20">
		<title lang="zh">Global Watch</title>
		<desc lang="zh"/>
	</programme>
</tv>
*/

// XmlTV tv tag
type XmlTV struct {
	Generator string   `xml:"generator-info-name,attr" json:"generatorInfoName"`
	Source    string   `xml:"source-info-name,attr" json:"sourceInfoName"`
	XMLName   xml.Name `xml:"tv" json:"-"`

	Channel []*XmlTvChannel `xml:"channel" json:"channel"`
	Program []*Program      `xml:"programme" json:"programme"`
}

// XmlTvChannel : channel info
type XmlTvChannel struct {
	ID          string        `xml:"id,attr" json:"id"`
	DisplayName []DisplayName `xml:"display-name" json:"displayName"`
}

// DisplayName 频道名
type DisplayName struct {
	Lang  string `xml:"lang,attr" json:"lang"`
	Value string `xml:",chardata" json:"value"`
}

// Program 节目
type Program struct {
	Channel string `xml:"channel,attr" json:"channel"`
	Start   string `xml:"start,attr" json:"start"`
	Stop    string `xml:"stop,attr" json:"stop"`

	Title []*Title `xml:"title" json:"title"`
	Desc  []*Desc  `xml:"desc" json:"desc"`
}

// Title 节目标题
type Title struct {
	Lang  string `xml:"lang,attr" json:"lang"`
	Value string `xml:",chardata" json:"value"`
}

// Desc : 节目描述
type Desc struct {
	Lang  string `xml:"lang,attr" json:"lang"`
	Value string `xml:",chardata" json:"value"`
}
