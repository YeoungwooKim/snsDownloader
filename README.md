# snsDownloader

#### 환경구성
> mongo db <br>
> fiber - [link](https://github.com/gofiber/fiber) <br>
> yt-dlp - [link](https://github.com/yt-dlp/yt-dlp) <br>
<hr/>

#### 주요 기능
* sns의 미디어 컨텐츠 링크로 영상 다운로드
	* 디폴트/ 해상도,,, etc 변경 가능.
* TODO - 유저별 컨텐츠 스크래핑 예약 기능
	* 서버가 바쁠땐 예약만 걸어놓고, 쉴때 처리하는 queue 처리 방안
	* 현재 서버에서 몇개의 일을 하고 있는지에 대해서 서버가 알아야됨

* TODO - 스크래핑 예약기능 - 등록/삭제
    * 히스토리 기능
* TODO - 미디어 컨텐츠 다운로드 실행시 client가 해당 작업 상태 확인 기능
	* 뭐,, 다운로드 몇 퍼센트 받았다,,, 뭐 트랜스코딩 몇 퍼센트 진행중이다,,, 등
	* 해당 태스크에 대해 유니크키가 존재해야 -> uuid로다가 ㄱㄱ
	* 진행 로직을 다음과 같이 해야 할거 같음

