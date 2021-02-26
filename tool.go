package util

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func computeTimeDiff(diff int64) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "now"
	case diff < 2:
		diff = 0
		diffStr = "1 second"
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d seconds", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 minute"
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d minutes", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 hour"
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d hours", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 day"
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d days", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 week"
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d weeks", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 month"
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d months", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 year"
	default:
		diffStr = fmt.Sprintf("%d years", diff/Year)
		diff = 0
	}
	return diff, diffStr
}

func TimeSincePro(then time.Time) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return "future"
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}

/*
Levenshtein 计算两个字符串之间的编辑距离，SimilarText 计算两个字符串的相似度
fmt.Println(Levenshtein("焦作市美佳百货有限公司", "焦作市美佳百货有限责任公司", 1, 1, 1))

	var percent float64
	fmt.Println(SimilarText("焦作市美佳百货有限公司", "焦作美佳百货有限责任公司", &percent))
	fmt.Println(percent)
*/

func Levenshtein(str1, str2 string, costIns, costRep, costDel int) int {
	var maxLen = 255
	l1 := len(str1)
	l2 := len(str2)
	if l1 == 0 {
		return l2 * costIns
	}
	if l2 == 0 {
		return l1 * costDel
	}
	if l1 > maxLen || l2 > maxLen {
		return -1
	}

	p1 := make([]int, l2+1)
	p2 := make([]int, l2+1)
	var c0, c1, c2 int
	var i1, i2 int
	for i2 := 0; i2 <= l2; i2++ {
		p1[i2] = i2 * costIns
	}
	for i1 = 0; i1 < l1; i1++ {
		p2[0] = p1[0] + costDel
		for i2 = 0; i2 < l2; i2++ {
			if str1[i1] == str2[i2] {
				c0 = p1[i2]
			} else {
				c0 = p1[i2] + costRep
			}
			c1 = p1[i2+1] + costDel
			if c1 < c0 {
				c0 = c1
			}
			c2 = p2[i2] + costIns
			if c2 < c0 {
				c0 = c2
			}
			p2[i2+1] = c0
		}
		tmp := p1
		p1 = p2
		p2 = tmp
	}
	c0 = p1[l2]

	return c0
}

// SimilarText similar_text()
func SimilarText(first, second string, percent *float64) int {
	var similarText func(string, string, int, int) int
	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText(first, second, l1, l2)
	if percent != nil {
		*percent = float64(sim*200) / float64(l1+l2)
	}
	return sim
}

// Subtract deals with subtraction of all types of number.
func Subtract(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt = true
	switch left := left.(type) {
	case int:
		rleft = int64(left)
	case int8:
		rleft = int64(left)
	case int16:
		rleft = int64(left)
	case int32:
		rleft = int64(left)
	case int64:
		rleft = left
	case float32:
		fleft = float64(left)
		isInt = false
	case float64:
		fleft = left
		isInt = false
	}

	switch right := right.(type) {
	case int:
		rright = int64(right)
	case int8:
		rright = int64(right)
	case int16:
		rright = int64(right)
	case int32:
		rright = int64(right)
	case int64:
		rright = right
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	if isInt {
		return rleft - rright
	} else {
		return fleft + float64(rleft) - (fright + float64(rright))
	}
}

// EllipsisString returns a truncated short string,
// it appends '...' in the end of the length of string is too large.
func EllipsisString(str string, length int) string {
	if len(str) < length {
		return str
	}
	return str[:length-3] + "..."
}

// TruncateString returns a truncated string with given limit,
// it returns input string if length is not reached limit.
func TruncateString(str string, limit int) string {
	if len(str) < limit {
		return str
	}
	return str[:limit]
}

func RangeCode() string {

	i := GenerateRangeNum(1000, 9999)
	str := strconv.Itoa(i)
	return str
}

func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
func LongToDate(timestamp int64) string {

	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05")
}
func IntToDate(timestamp int64, fm string) string {

	if fm == "" {
		fm = "20060102"
	}
	tm := time.Unix(timestamp, 0)
	return tm.Format(fm)
}

func OrderNo(now int64) string {

	dt := time.Now().Format("060102")
	//tmp := now
	dt1 := time.Now().Format("2006-01-02")
	tm, _ := time.Parse("2006-01-02 15:04:05", dt1+" 00:00:00")
	tmp := tm.UnixNano() - 28800000000000
	//fmt.Println("now:", now)
	//fmt.Println("tmp:", tmp)
	tmp = now - tmp

	payno := dt + strconv.FormatInt(tmp/10000, 10)
	//1805294200873517
	//9007199254740992 js最大精度
	return payno
}

//184.68000000000006 -> 184.68
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

//0.01800000000000 -> 1.80%
func Rate(src string) string {

	f64, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return ""
	}
	f64 = f64 * 100
	value := fmt.Sprintf("%.2f", f64)
	return value + "%"
}

func F64ToInt64(f float64) int64 {

	s := fmt.Sprintf("%0.0f", f)
	v, _ := strconv.ParseFloat(s, 64)
	return int64(v)
}

//四舍六入五成双
func F2ncy(f float64) float64 {

	//f = f + 0.005
	s := fmt.Sprintf("%0.2f", f)
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func Long(f float64) int64 {

	return int64(math.Floor(f + 0/5))
}

func Ncy2Int(f64 float64) int64 {

	s := fmt.Sprintf("%0.0f", f64*100)
	i64, _ := strconv.ParseInt(s, 10, 64)
	return int64(i64)
}
func Int2ncy(i int64) float64 {

	var f float64 = float64(i)
	f = f / 100
	s := fmt.Sprintf("%0.2f", f)
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
func Int2NcyStr(i int64) string {

	if i == 0 {
		return "0.00"
	}
	var f float64 = float64(i)
	f = f / 100
	s := fmt.Sprintf("%0.2f", f)
	return s
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substring(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
func FloatAbs(n float64) float64 {
	if n < 0 {
		return -1 * n
	}
	return n
}
func IntAbs(n int64) int64 {
	y := n >> 63       // y ← x >> 63
	return (n ^ y) - y // (x ⨁ y) - y
}
