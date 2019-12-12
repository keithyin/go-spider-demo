package main

import (
	"fmt"
	"math/rand"
)

// var urlList = []string{
// 	"https://www.weibo.com/3093406867/profile?rightmod=1&wvr=6&mod=personinfo&is_all=1",
// 	"https://www.weibo.com/p/aj/v6/mblog/mbloglist?ajwvr=6&domain=100505&rightmod=1&wvr=6&mod=personinfo&is_all=1&pagebar=0&pl_name=Pl_Official_MyProfileFeed__20&id=1005053093406867&script_uri=/3093406867/profile&feed_type=0&page=%d&pre_page=1&domain_op=100505&__rnd=%d",
// }
var urlList = []string{
	"https://www.weibo.com/5358290784/profile?rightmod=1&wvr=6&mod=personnumber&is_all=1",
	"https://www.weibo.com/p/aj/v6/mblog/mbloglist?ajwvr=6&domain=100505&rightmod=1&wvr=6&mod=personnumber&is_all=1&pagebar=0&pl_name=Pl_Official_MyProfileFeed__20&id=1005055358290784&script_uri=/5358290784/profile&feed_type=0&page=1&pre_page=%d&domain_op=100505&__rnd=%d",
}

func GetCookie() string {
	cookie := "Apache=9853586515057.533.1541596560816; SINAGLOBAL=9853586515057.533.1541596560816; ULV=1541596560819:1:1:1:9853586515057.533.1541596560816:; _s_tentry=passport.weibo.com; login_sid_t=a7e9289cf6693a372395919a84533a1a; cross_origin_proto=SSL; YF-Ugrow-G0=8751d9166f7676afdce9885c6d31cd61; WBStorage=e8781eb7dee3fd7f|undefined; UOR=,,cn.bing.com; YF-V5-G0=a2489c19ecf98bbe86a7bf6f0edcb071; wb_view_log=1280*7201.5; SUB=_2A2525pZQDeRhGeNN7loT-S7LwziIHXVVlYCYrDV8PUNbmtBeLXijkW9NSbFGNmrhPSkLadaC41zjnmHGBSJGhUBo; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WFB.rWULTh38K4vJz3RGNyH5JpX5KzhUgL.Fo-0SKnE1K5N1hB2dJLoI0qLxKqL1hnL1K5LxK.L1heLB.BLxKML12eLB-zLxK-L1KzLB-2LxK.LBKzL1hqLxK-L1h-LB.zt; SUHB=0RyqK6VRYZwvSm; ALF=1573132672; SSOLoginState=1541596673; wvr=6; YF-Page-G0=854ebb7f403eecfa60ed1f0e977c6825; wb_view_log_5358290784=1280*7201.5"
	return cookie
}

func GenerateUrl(page_index int) (string, int) {
	if page_index == 0 {
		return urlList[page_index], page_index
	} else {
		return fmt.Sprintf(urlList[1], page_index, rand.Uint64()), page_index
	}
}
