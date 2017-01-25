package mal

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://myanimelist.net"

type User struct {
	Id       int    `xml:"id"`
	Username string `xml:"username"`
}

type MAL struct {
	Username string
	Password string
}

func NewMAL(username, password string) *MAL {
	return &MAL{Username: username, Password: password}
}

type AnimeSearch struct {
	Anime []Anime `xml:"entry"`
}

type Anime struct {
	Id       int     `xml:"id"`
	Title    string  `xml:"title"`
	English  string  `xml:"english"`
	Synonyms string  `xml:"synonyms"`
	Episodes int     `xml:"episodes"`
	Type     string  `xml:"type"`
	Status   string  `xml:"status"`
	Start    MALDate `xml:"start_date"`
	End      MALDate `xml:"end_date"`
	Synopsis string  `xml:"synopsis"`
	ImageURL string  `xml:"image"`
}

type AnimePayload struct {
	/*
		episode. int
		status. int OR string. 1/watching, 2/completed, 3/onhold, 4/dropped, 6/plantowatch
		score. int
		storage_type. int (will be updated to accomodate strings soon)
		storage_value. float
		times_rewatched. int
		rewatch_value. int
		date_start. date. mmddyyyy
		date_finish. date. mmddyyyy
		priority. int
		enable_discussion. int. 1=enable, 0=disable
		enable_rewatching. int. 1=enable, 0=disable
		comments. string
		fansub_group. string
		tags. string. tags separated by commas
	*/

	Episode          int       `xml:"episode"`
	Status           string    `xml:"status"`
	Score            int       `xml:"score"`
	StorageType      int       `xml:"storage_type"`
	StorageValue     int       `xml:"storage_value"`
	TimesRewatched   int       `xml:"times_rewatched"`
	RewatchValue     int       `xml:"rewatch_value"`
	Start            time.Time `xml:"date_start"`
	End              time.Time `xml:"date_end"`
	Priority         int       `xml:"priority"`
	EnableDiscussion bool      `xml:"enable_discussion"`
	EnableRewatching bool      `xml:"enable_rewatching"`
	Comments         string    `xml:"comments"`
	FansubGroup      string    `xml:"fansub_group"`
	Tags             []string  `xml:"tags"`
}

func (ap AnimePayload) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "entry"}
	e.EncodeToken(start)
	e.EncodeElement(ap.Episode, xml.StartElement{Name: xml.Name{Local: "episode"}})
	e.EncodeElement(ap.Status, xml.StartElement{Name: xml.Name{Local: "status"}})
	e.EncodeElement(ap.Score, xml.StartElement{Name: xml.Name{Local: "score"}})
	e.EncodeElement(ap.StorageType, xml.StartElement{Name: xml.Name{Local: "storage_type"}})
	e.EncodeElement(ap.StorageValue, xml.StartElement{Name: xml.Name{Local: "storage_value"}})
	e.EncodeElement(ap.TimesRewatched, xml.StartElement{Name: xml.Name{Local: "times_rewatched"}})
	e.EncodeElement(ap.RewatchValue, xml.StartElement{Name: xml.Name{Local: "rewatch_value"}})
	e.EncodeElement(ap.Priority, xml.StartElement{Name: xml.Name{Local: "priority"}})
	e.EncodeElement(ap.Comments, xml.StartElement{Name: xml.Name{Local: "comments"}})
	e.EncodeElement(ap.FansubGroup, xml.StartElement{Name: xml.Name{Local: "FansubGroup"}})
	e.EncodeElement(strings.Join(ap.Tags, ","), xml.StartElement{Name: xml.Name{Local: "tags"}})

	dateFormat := "01022006"
	e.EncodeElement(ap.Start.Format(dateFormat), xml.StartElement{Name: xml.Name{Local: "date_start"}})
	e.EncodeElement(ap.End.Format(dateFormat), xml.StartElement{Name: xml.Name{Local: "date_end"}})

	enableDiscussion := 0
	if ap.EnableDiscussion {
		enableDiscussion = 1
	}
	e.EncodeElement(enableDiscussion, xml.StartElement{Name: xml.Name{Local: "enable_discussion"}})

	enableRewatching := 0
	if ap.EnableRewatching {
		enableRewatching = 1
	}
	e.EncodeElement(enableRewatching, xml.StartElement{Name: xml.Name{Local: "enable_rewatching"}})

	e.EncodeToken(xml.EndElement{Name: start.Name})

	return nil
}

type MyAnimeList struct {
	UserID    int                `xml:"myinfo>user_id"`
	Username  string             `xml:"myinfo>user_name"`
	Watching  int                `xml:"myinfo>user_watching"`
	Completed int                `xml:"myinfo>user_completed"`
	Hold      int                `xml:"myinfo>user_onhold"`
	Dropped   int                `xml:"myinfo>user_dropped"`
	Plan      int                `xml:"myinfo>user_plantowatch"`
	Days      float64            `xml:"myinfo>user_days_spent_watching"`
	Anime     []MyAnimeListAnime `xml:"anime"`
}

type MyAnimeStatus int
type MyAnimeRewatch int
type MyAnimePriority int

func (m *MyAnimeRewatch) Unmarshal(d *xml.Decoder, start xml.StartElement) error {
	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}

	*m = MyAnimeRewatch(i)
	return nil
}

func (m *MyAnimeRewatch) Name() string {
	switch *m {
	case 1:
		return "Very Low"
	case 2:
		return "Low"
	case 3:
		return "Medium"
	case 4:
		return "High"
	case 5:
		return "Very High"
	default:
		return "None"
	}
}

func (m *MyAnimeRewatch) Val() int {
	return int(*m)
}

func (m MyAnimeRewatch) String() string {
	return m.Name()
}

func (m *MyAnimeStatus) Name() string {
	switch *m {
	case 1:
		return "Watching"
	case 2:
		return "Completed"
	case 3:
		return "On Hold"
	case 4:
		return "Dropped"
	case 6:
		return "Plan to Watch"
	default:
		return "Unknown"
	}
}

func (m MyAnimeStatus) String() string {
	return m.Name()
}

func (m *MyAnimeStatus) Val() int {
	return int(*m)
}

func (m *MyAnimeStatus) Unmarshal(d *xml.Decoder, start xml.StartElement) error {
	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}

	*m = MyAnimeStatus(i)
	return nil
}

type MyAnimeListAnime struct {
	AnimeDBID           int            `xml:"series_animedb_id"`
	Title               string         `xml:"series_title"`
	Synonyms            string         `xml:"series_synonyms"`
	Type                int            `xml:"series_type"`
	Episodes            int            `xml:"series_episodes"`
	Status              int            `xml:"series_status"`
	Start               MALDate        `xml:"series_start"`
	End                 MALDate        `xml:"series_end"`
	ImageURL            string         `xml:"series_image"`
	MyID                int            `xml:"my_id"`
	MyWatched           int            `xml:"my_watched_episodes"`
	MyStart             MALDate        `xml:"my_start_date"`
	MyFinish            MALDate        `xml:"my_finish_date"`
	MyScore             int            `xml:"my_score"`
	MyStatus            MyAnimeStatus  `xml:"my_status"`
	MyRewatch           MyAnimeRewatch `xml:"my_rewatching"`
	MyRewatchingEpisode int            `xml:"my_rewatching_ep"`
	MyUpdated           MALDateTime    `xml:"my_last_updated"`
	MyTags              string         `xml:"my_tags"`
}

type MALDateTime struct {
	time.Time
}

func (mdt *MALDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	mdt.Time = time.Unix(i, 0)
	return nil
}

type MALDate struct {
	time.Time
}

func (m *MALDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	m.Time, _ = time.Parse("2006-01-02", s)
	return nil
}

func (mal *MAL) Verify() (*User, error) {
	req, err := http.NewRequest("GET", baseURL+"/api/account/verify_credentials.xml", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(mal.Username, mal.Password)

	client := &http.Client{Timeout: time.Second * 30}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 204 {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	user := User{}
	err = xml.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (mal *MAL) Search(q string) (*AnimeSearch, error) {
	req, err := http.NewRequest("GET", baseURL+"/api/anime/search.xml", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(mal.Username, mal.Password)

	query := req.URL.Query()
	query.Add("q", q)

	req.URL.RawQuery = query.Encode()

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 204 {
		return nil, errors.New("204: No content")
	}

	defer res.Body.Close()

	result := AnimeSearch{}
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (mal *MAL) MyAnime() (*MyAnimeList, error) {
	req, err := http.NewRequest("GET", baseURL+"/malappinfo.php", nil)
	if err != nil {
		return nil, err
	}
	// ???: This resource does not need base authentication.
	req.SetBasicAuth(mal.Username, mal.Password)

	query := req.URL.Query()
	query.Add("u", mal.Username)
	query.Add("status", "all")
	query.Add("type", "anime")

	req.URL.RawQuery = query.Encode()

	client := http.Client{Timeout: time.Second * 30}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	myAnimeList := MyAnimeList{}
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(&myAnimeList)
	if err != nil {
		return nil, err
	}

	return &myAnimeList, nil
}

func (mal *MAL) Add(id int, ap AnimePayload) error {
	b, err := xml.Marshal(ap)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("data", string(b))

	url := fmt.Sprintf("%s/api/animelist/add/%d.xml", baseURL, id)
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(mal.Username, mal.Password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Second * 30}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		return errors.New(res.Status)
	}

	return nil
}

func (mal *MAL) Update(id int, ap AnimePayload) error {
	b, err := xml.Marshal(ap)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("data", string(b))

	url := fmt.Sprintf("%s/api/animelist/update/%d.xml", baseURL, id)
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(mal.Username, mal.Password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Second * 30}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		return errors.New(res.Status)
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	body := string(bytes)
	if body != "Updated" {
		return errors.New("Unexpected response body: " + body)
	}

	return nil
}

func (mal *MAL) Delete(id int) error {
	url := fmt.Sprintf("%s/api/animelist/delete/%d.xml", baseURL, id)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(mal.Username, mal.Password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Second * 30}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		return errors.New(res.Status)
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	body := string(bytes)
	if body != "Deleted" {
		return errors.New("Unexpected response body: " + body)
	}

	return nil

}
