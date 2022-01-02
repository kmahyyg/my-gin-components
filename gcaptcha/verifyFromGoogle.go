package gcaptcha

import (
	"encoding/json"
	"errors"
	"github.com/kmahyyg/my-gin-components/common-conf"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrVerifierNotBuilt = errors.New("verifier is not built yet")
	ErrNotAHuman = errors.New("parse google response success but not meet the requirement of human threshold")
	ErrInvalidThreshold = errors.New("invalid threshold, should >0.0 and <=1.0")
	ErrCannotVerify = errors.New("verification to Google Server Failed. Check site secret is corresponding to your domain or network connection")
)

type gRecaptchaResponse struct {
	Success bool `json:"success"`
	Timestamp int `json:"challenge_ts"`
	Hostname string `json:"hostname"`
	Score float32 `json:"score"`
}

type gRecaptchaVerifier struct {
	threshold float32
	inChinaMainland bool
	siteSecret string
	verifyEndpoint string
	domainName string
}

// GCaptchaVerifierFactory build a singleton of captcha verifier
type GCaptchaVerifierFactory struct {
	isBuilt bool
	verifier *gRecaptchaVerifier
}

func (gcvf *GCaptchaVerifierFactory) GetVerifier() (*gRecaptchaVerifier,error) {
	if gcvf.isBuilt && gcvf.verifier != nil {return gcvf.verifier, nil}
	return nil, ErrVerifierNotBuilt
}

// BuildVerifier return a singleton of gRecaptchaVerifier
// @param thres: threshold, smaller than this value will be interpreted as Bot and get banned.
// @param inChina: China Mainland should set this value to true.
// @param siteSecret: Google ReCaptcha Verifier
func (gcvf *GCaptchaVerifierFactory) BuildVerifier(conf common_conf.ReCaptchaConfig) {
	if conf.Threshold > float32(1.0) || conf.Threshold <= float32(0.0) {
		panic(ErrInvalidThreshold)
	}
	var siteEndP = "https://www.google.com/recaptcha/api/siteverify"
	if conf.InChina {
		siteEndP = "https://recaptcha.net/recaptcha/api/siteverify"
	}
	gcvf.verifier = &gRecaptchaVerifier{
		threshold:       conf.Threshold,
		inChinaMainland: conf.InChina,
		siteSecret:      conf.Secret,
		verifyEndpoint:  siteEndP,
		domainName: 	 conf.DomainName,
	}
	gcvf.isBuilt = true
}

// ResetVerifier works if any config went wrong, this function will set gRecaptchaVerifier to nil in Factory.
// after this, call GetVerifier to create one new verifier.
func (gcvf *GCaptchaVerifierFactory) ResetVerifier(){
	gcvf.isBuilt = false
	gcvf.verifier = nil
	return
}

// Verify method is implemented to gRecaptchaVerifier sent
// @param gcaptresp: string, g-captcha-response data
// @return verified: bool, successfully verified, return true
// @return error: error, if url request error occurred or Google refused us will return not nil.
func (gcapt *gRecaptchaVerifier) Verify(gcaptresp string) (bool,error) {
	resp, err := http.PostForm(gcapt.verifyEndpoint, url.Values{
		"secret": {gcapt.siteSecret},
		"response": {gcaptresp},
	})
	if err != nil {return false,err}
	if resp.StatusCode != 200 {return false, ErrCannotVerify}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	respStru := &gRecaptchaResponse{}
	err = json.Unmarshal(body, respStru)
	if err != nil {return false, ErrCannotVerify}
	// loosely check if domain name is corresponding
	if respStru.Success && strings.Contains(respStru.Hostname, gcapt.domainName) && respStru.Score > gcapt.threshold {
		return true, nil
	}
	return false, ErrNotAHuman
}
