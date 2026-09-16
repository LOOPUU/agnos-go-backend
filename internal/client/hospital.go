package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/LOOPUU/agnos-hospital-middleware/internal/model"
)

type HospitalClient struct { HTTP *http.Client }
func New(timeout time.Duration) *HospitalClient { return &HospitalClient{HTTP:&http.Client{Timeout:timeout}} }
func (c *HospitalClient) Find(ctx context.Context, baseURL,id string)(*model.Patient,error){
	req,err:=http.NewRequestWithContext(ctx,http.MethodGet,strings.TrimRight(baseURL,"/")+"/patient/search/"+url.PathEscape(id),nil);if err!=nil{return nil,err}
	req.Header.Set("Accept","application/json");resp,err:=c.HTTP.Do(req);if err!=nil{return nil,fmt.Errorf("hospital API: %w",err)};defer resp.Body.Close()
	if resp.StatusCode==http.StatusNotFound{return nil,nil};if resp.StatusCode!=http.StatusOK{return nil,fmt.Errorf("hospital API returned %d",resp.StatusCode)}
	var p model.Patient;if err=json.NewDecoder(resp.Body).Decode(&p);err!=nil{return nil,fmt.Errorf("invalid hospital response: %w",err)};if p.PatientHN==""{return nil,fmt.Errorf("hospital response missing patient_hn")};return &p,nil
}

