// Package general - contains general functions
package general

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cr "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"olx-clone/functions/logger"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var log = logger.Log

// ValidateStruct validates the struct and return error if they occur
func ValidateStruct(s interface{}) error {
	v := validator.New()
	return v.Struct(s)
}

// ValidUserName return bool if the username if valid, to prevent SQL injection
func ValidUserName(username string) bool {
	list := [4]string{"'", "--", "void", "null"}
	for _, v := range list {
		if strings.Contains(username, v) {
			return false
		}
	}
	return true
}

// StartTime returns current time
func StartTime() time.Time {
	return time.Now()
}

// TimeDifference returns elapsed time
func TimeDifference(startTime time.Time) time.Duration {
	return time.Since(startTime)
}

// LogTimeDifference logs the elapsed time
func LogTimeDifference(startTime time.Time) {
	elapsedTime := TimeDifference(startTime)
	log.Println("Time taken: ", elapsedTime)
}

// GetUUID generates and returns a new UUID v4 string
func GetUUID() string {
	return uuid.New().String()
}

// GenerateRandomBytes generates and returns random bytes of specified size n
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := cr.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString generates and returns a random alphanumeric string with a start_string (with mixed case)
// of size n
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

// InArrStr is func that checks for val in Array of type of string
func InArrStr(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// InArrInt is func that checks for val in Array of type of int
func InArrInt(val int, arr []int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

var spaceRegex = regexp.MustCompile(`\s+`)

// RemoveExtraSpaces removes excess spaces from a string
func RemoveExtraSpaces(input string) string {
	return strings.TrimSpace(spaceRegex.ReplaceAllString(input, " "))
}

// RemoveAllSpaces removes all spaces from a string
func RemoveAllSpaces(input string) string {
	return strings.TrimSpace(spaceRegex.ReplaceAllString(input, ""))
}

// RemoveEmptyElements removes empty values from list
func RemoveEmptyElements(s []string) []string {
	var r []string
	for _, str := range s {
		str = RemoveExtraSpaces(str)
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// RemoveDuplicatesAndEmpty removes duplicate strings from list
func RemoveDuplicatesAndEmpty(s []string) []string {
	var r []string
	uniqueValues := make(map[string]bool)
	for _, str := range s {
		str = RemoveExtraSpaces(str)
		if str != "" {
			found := uniqueValues[str]
			if !found {
				uniqueValues[str] = true
				r = append(r, str)
			}
		}
	}
	return r
}

var ValidLanguageNameRegex = regexp.MustCompile("^[a-zA-Z ]{3,70}$")

var cleanAddressRegex = regexp.MustCompile(`[^a-zA-Z\-\,\.\s\@\#\(\)\/0-9]|\n+`)

// CleanAddressLine removes unallowed characters from address
func CleanAddressLine(addressLine string) string {
	return RemoveExtraSpaces(cleanAddressRegex.ReplaceAllString(addressLine, ""))
}

var cleanColumnRegex = regexp.MustCompile(`[^a-zA-Z\-\_0-9]+`)

// CleanColumn removes unallowed characters from address
func CleanColumn(column string) string {
	return RemoveExtraSpaces(cleanColumnRegex.ReplaceAllString(column, ""))
}

var keepOnlyAlphaSpaceRegex = regexp.MustCompile(`[^a-zA-Z\s]+`)

// KeepOnlyAlphaSpace removes everything except alphabets and space
func KeepOnlyAlphaSpace(text string) string {
	return RemoveExtraSpaces(keepOnlyAlphaSpaceRegex.ReplaceAllString(text, ""))
}

var alphabetNumericExcludeRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// GetOnlyAlphaNumUpper removes everything except alpha numeric characters and turn to upper case
func GetOnlyAlphaNumUpper(input string) string {
	return strings.ToUpper(alphabetNumericExcludeRegex.ReplaceAllString(input, ""))
}

var alphabetNumSpaceExcludeRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)

// GetOnlyAlphaNumSpace removes everything except alpha numeric characters and space
// it also removes extra spaces, and keeps case intact
func GetOnlyAlphaNumSpace(input string) string {
	return alphabetNumSpaceExcludeRegex.ReplaceAllString(RemoveExtraSpaces(input), "")
}

var alphabetSpaceExcludeRegex = regexp.MustCompile(`[^a-zA-Z\s]+`)

// GetOnlyAlphaSpace removes everything except alpha characters and space
// it also removes extra spaces, and keeps case intact
func GetOnlyAlphaSpace(input string) string {
	return strings.TrimSpace(alphabetSpaceExcludeRegex.ReplaceAllString(RemoveExtraSpaces(input), ""))
}

var countryCodeRegex = regexp.MustCompile(`^\+91`)

// RemoveCountryCode removes the country code and spaces from input mobile string
func RemoveCountryCode(mobile string) string {
	return countryCodeRegex.ReplaceAllString(RemoveExtraSpaces(mobile), "")
}

// GetStringFromTemplate constructs a string by replacing placeholders using a map and returns it
func GetStringFromTemplate(templateString string, data map[string]interface{}) string {
	t := template.Must(template.New("temptemplate").Parse(templateString))
	builder := &strings.Builder{}
	if err := t.Execute(builder, data); err != nil {
		return ""
	}
	return builder.String()
}

// IsNumber returns true if all chars in s are digits
func IsNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// GetTimeStampString returns the current time in YYYY-MM-DD HH:MM:SS.MMMMMM format
func GetTimeStampString() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}

// GetTimeStampPair returns the current time object as well as YYYY-MM-DD HH:MM:SS.MMMMMM format
func GetTimeStampPair() (time.Time, string) {
	date := time.Now()
	return date, date.Format("2006-01-02 15:04:05.000000")
}

var pinCodeRegex = regexp.MustCompile(`^[1-9][0-9]{5}$`)

// ValidatePincode validates input string for being a valid pincode
func ValidatePincode(pincode string) bool {
	return pinCodeRegex.MatchString(pincode)
}

var gstinRegex = regexp.MustCompile(`^([0-2][0-9]|[3][0-7])[A-Z]{3}[ABCFGHLJPTK][A-Z]\d{4}[A-Z][A-Z0-9][Z][A-Z0-9]$`)

// ValidateGSTIN validates input capitalised string for being a valid GSTIN
func ValidateGSTIN(gstin string) bool {
	return gstinRegex.MatchString(gstin)
}

var panRegex = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)

// ValidatePAN validates input capitalised string for being a valid PAN
func ValidatePAN(pan string) bool {
	return panRegex.MatchString(pan)
}

// ValidatePersonalPAN validates input capitalised string for being a valid PAN and if its personal
func ValidatePersonalPAN(pan string) bool {
	isPAN := panRegex.MatchString(pan)
	if isPAN {
		return pan[3:4] == "P"
	}
	return false
}

// ValidateNonPersonalPAN validates input capitalised string for being a valid PAN
func ValidateNonPersonalPAN(pan string) bool {
	isPAN := panRegex.MatchString(pan)
	if isPAN {
		return pan[3:4] != "P"
	}
	return false
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// ValidateEmail validates an email
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidateURL validates input string for being a valid URL
func ValidateURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// GetFirstSurName returns first and last name
func GetFirstSurName(name string) (string, string) {
	name = RemoveExtraSpaces(name)
	var firstName string
	var surName string
	spaceIndex := strings.LastIndex(name, " ")
	if spaceIndex == -1 {
		firstName = name
		surName = ""
	} else {
		firstName = name[:spaceIndex]
		surName = name[spaceIndex+1:]
	}
	return firstName, surName
}

// AmountInWords returns the amount in words
func AmountInWords(amount float64) string {
	return spellFloat(amount)
}

func spellFloat(number float64) string {
	i := int(number)
	d := int((number - float64(i)) * 100)
	if d > 0 {
		return fmt.Sprintf("%s and %s paisa", spell(i), spell(d))
	}
	return spell(i)
}

func spell(n int) string {
	to19 := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve",
		"thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens := []string{"twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	if n == 0 {
		return ""
	}
	if n < 20 {
		return to19[n-1]
	}
	if n < 100 {
		if n%10 == 0 {
			return tens[n/10-2]
		}
		return tens[n/10-2] + "-" + spell(n%10)
	}
	if n < 1000 {
		return to19[n/100-1] + " hundred " + spell(n%100)
	}
	if n < 100000 {
		return spell(n/1000) + " thousand " + spell(n%1000)
	}
	if n < 10000000 {
		anchor := n / 100000
		if anchor > 1 {
			return fmt.Sprintf("%s lakhs %s", spell(anchor), spell(n%100000))
		}
		return fmt.Sprintf("%s lakh %s", spell(anchor), spell(n%100000))
	}
	if n < 1000000000 {
		anchor := n / 10000000
		if anchor > 1 {
			return fmt.Sprintf("%s crores %s", spell(anchor), spell(n%10000000))
		}
		return fmt.Sprintf("%s crore %s", spell(anchor), spell(n%10000000))
	}
	return ""
}

// IsStringNumeric returns if string contains only numbers
func IsStringNumeric(inputStr string) bool {
	isNotDigit := func(c rune) bool { return c < '0' || c > '9' }
	return strings.IndexFunc(inputStr, isNotDigit) == -1
}

// AESEncrypt does AES Encryption and return base64
func AESEncrypt(src, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Errorln("key error1", err)
	}
	if src == "" {
		log.Errorln("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	// Normal base64 encoding encryption is different from urlsafe base64
	return base64.StdEncoding.EncodeToString(crypted)
}

// AESDecrypt does decryption of base64 encoded AES encrypted data an returns plain string
func AESDecrypt(b64src string, key []byte) string {
	src, err := base64.StdEncoding.DecodeString(b64src)
	if err != nil {
		log.Errorln("error decoding base64 content")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Errorln("err is:", err)
	}
	blockMode := NewECBDecrypter(block)
	crypted := []byte(src)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)

	return string(origData)
}

func AESCBCPKCS5Encryption(message, key string) (string, error) {
	if message == "" || key == "" {
		return "", errors.New("message or key cannot be empty")
	}
	messageBytes := []byte(message)
	keybytes := []byte(key)
	vector := make([]byte, 16)
	encodedMessage := make([]byte, len(messageBytes)+16)
	copy(encodedMessage[0:], vector[0:])
	copy(encodedMessage[16:], messageBytes[0:])
	c, err := aes.NewCipher(keybytes)
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	ecb := cipher.NewCBCEncrypter(c, make([]byte, 16))
	content := PKCS5Padding(encodedMessage, c.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	res := base64.StdEncoding.EncodeToString(crypted)
	return res, nil
}

func AESCBCPKCS5Decryption(encryptedMessage, key string) (string, error) {
	if encryptedMessage == "" || key == "" {
		return "", errors.New("encryptedMessage or key cannot be empty")
	}
	decodedEncryptedMessage, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	keyBytes := []byte(key)
	vector := make([]byte, 16)
	decodedEncrypted := make([]byte, len(decodedEncryptedMessage)-16)
	copy(decodedEncrypted[0:], decodedEncryptedMessage[16:])
	copy(vector[0:], decodedEncryptedMessage[0:len(vector)])
	c, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	ecb := cipher.NewCBCDecrypter(c, vector)
	decryptedMessage := make([]byte, len(decodedEncrypted))
	ecb.CryptBlocks(decryptedMessage, decodedEncrypted)
	return string(PKCS5UnPadding(decryptedMessage)), nil
}

func GenerateHashedKey(password, key string) []byte {
	// hash text key using sha256
	sum := sha256.Sum256([]byte(key))
	// convert from [32]byte to []byte
	salt := sum[:32]
	// generete secret key
	secretKey := pbkdf2.Key([]byte(password), salt, 65536, 32, sha256.New)
	return secretKey
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// remove the last byte unpadding times
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// FileFromURLtoBase64 retrieves file using url and converts to base64
func FileFromURLtoBase64(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Errorln(err)
		return ""
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func DeleteFileByPath(path string) error {
	return os.Remove(path)
}

// IdentReader is used while xml parsing
func IdentReader(encoding string, input io.Reader) (io.Reader, error) {
	return input, nil
}

// FormatCurrency returns a formatted string with commas for given amount
func FormatCurrency(amount float64, showDecimal bool) string {
	p := message.NewPrinter(language.MustParse("en-IN"))
	if showDecimal {
		return p.Sprintf("%.2f", amount)
	}
	return p.Sprintf("%.0f", amount)
}

// ValidateAccountNumber checks for valid account number
func ValidateAccountNumber(accountNumber string) bool {
	if !IsStringNumeric(accountNumber) {
		return false
	}
	return len(accountNumber) >= 9 && len(accountNumber) <= 18
}

// CalculateMedian returns median of numbers
func CalculateMedian(n []float64) float64 {
	if len(n) == 0 {
		return 0
	}
	if len(n) == 1 {
		return n[0]
	}
	sort.Float64s(n) // sort the numbers
	mNumber := len(n) / 2
	if len(n)%2 != 0 { // odd
		return n[mNumber]
	}

	return (n[mNumber-1] + n[mNumber]) / 2
}

var removeNonUTF = func(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}

// RemoveNonUTF8Strings removes strings that isn't UTF-8 encoded
func RemoveNonUTF8Strings(string string) string {
	return strings.Map(removeNonUTF, string)
}

// RemoveNonUTF8Bytes removes bytes that isn't UTF-8 encoded
func RemoveNonUTF8Bytes(data []byte) []byte {
	return bytes.Map(removeNonUTF, data)
}

// IsFreeMail checks whether a given email is from a free service or not
func IsFreeMail(email string) bool {
	lowerEmail := strings.ToLower(email)
	isFreeEmail, _ := regexp.MatchString(`@(live|hotmail|outlook|aol|yahoo|rocketmail|gmail|gmx|mail|inbox|icloud|aim|zoho|yandex|rediffmail)\.`, lowerEmail)
	// also check for indian university emails
	if isFreeEmail || strings.HasSuffix(lowerEmail, ".ac.in") || strings.HasSuffix(lowerEmail, ".edu") || strings.HasSuffix(lowerEmail, ".edu.in") || strings.HasSuffix(lowerEmail, ".org") {
		return true
	}
	return false
}

// CheckWithinNDays takes input string in YYYY-MM-DD format
func CheckWithinNDays(inputDateStr string, n int) bool {
	inputDate, _ := time.Parse("2006-01-02", inputDateStr)
	return int(time.Since(inputDate).Hours()/24.0) <= n
}

func ConvertTimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func ConvertTimeToTimestamp(time time.Time) int64 {
	return time.UnixNano() / 1_000_000_000
}

func GetDifferenceInMonths(a, b time.Time) int {
	months := 0
	month := b.Month()
	for b.Before(a) {
		b = b.Add(time.Hour * 24)
		nextMonth := b.Month()
		if nextMonth != month {
			months++
		}
		month = nextMonth
	}

	return months
}

// GetInitial returns two letter initials from a word
func GetInitial(word string) string {
	word = KeepOnlyAlphaSpace(word)

	words := []rune(word)
	if len(words) < 2 {
		return ""
	}
	firstLetter := words[0]
	secondLetter := words[1]
	for i := 2; i < len(words); i++ {
		if fmt.Sprintf("%c", words[i]) == " " {
			secondLetter = words[i+1]
			break
		}
	}
	return strings.ToUpper(fmt.Sprintf("%c%c", firstLetter, secondLetter))
}

func IsLeapYear(year int) bool {
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var ifscRegex = regexp.MustCompile(`[A-Z]{4}[A-Z0-9]{7}$`)

func ValidateIFSC(ifsc string) bool {
	return ifscRegex.MatchString(strings.ToUpper(ifsc))
}

func ValidateUUID(inputStr string) bool {
	_, err := uuid.Parse(inputStr)
	return err == nil
}

// CheckSupportedFormats checks for supported formats and returns file object, extension, error http code and error message if any
func CheckSupportedFormats(r *http.Request, extensions, contentTypes []string) (multipart.File, string, int, string) {
	// check for file types using file name
	file, header, _ := r.FormFile("file")
	if header == nil {
		// logger.WithRequest(r).Errorln("header is null")
		return nil, "", http.StatusBadRequest, "couldn't determine the file headers, please try again"
	}
	fileNameArr := strings.Split(header.Filename, ".")
	if len(fileNameArr) < 2 {
		// no extension case
		return nil, "", http.StatusBadRequest, "no file extension found"
	}
	extension := strings.ToLower(fileNameArr[len(fileNameArr)-1]) // change case to lower for easy checks
	if !InArrStr(extension, extensions) {                         // []string{"pdf", "jpeg", "png", "jpg"}
		return nil, "", http.StatusBadRequest, "file type not supported"
	}

	// now check using content type
	buff := make([]byte, 512) // why 512 bytes, see http://golang.org/pkg/net/http/#DetectContentType
	_, err := file.Read(buff)
	if err != nil {
		// logger.WithRequest(r).Errorln(err)
		return nil, "", http.StatusBadRequest, "couldn't read the file, please try again"
	}
	contentType := strings.ToLower(http.DetectContentType(buff))
	if !InArrStr(contentType, contentTypes) { // []string{"image/jpeg", "image/jpg", "image/png", "application/pdf"}
		return nil, "", http.StatusBadRequest, "file type not supported"
	}

	// reset the read pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		// logger.WithRequest(r).Errorln(err)
	}

	return file, extension, 0, ""
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func GetFilePathFromMultipart(file multipart.File) (string, error) {
	fileName := GetUUID()
	path := "/tmp/" + fileName
	defer file.Close()
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return path, nil
}

func ValidateName(name string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9\\s.]{2,99}$", name)
	return match
}

func ValidateReferenceID(referenceID string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]{2,99}$", referenceID)
	return match
}
