package sync

import log "github.com/sirupsen/logrus"

type SyncWorker struct {
}

func NewWorker() {

}

func (s *SyncWorker) Sync() {
	// TODO: caching download list, lol
	isBackingOff, err := dlReq.IsBackingOff()
	if err != nil {
		return err
	}

	// refresh cache if backoff period is up
	if !isBackingOff {
		log.Infof("Backoff period expired for download request %s, syncing all", dlReq.Id)
		itemsAdded, err := d.syncDownloadList(dlReq)
		if err != nil {
			return err
		}

		if itemsAdded {
			err = dlReq.ReportSyncHit()
			if err != nil {
				return err
			}
		} else {
			err = dlReq.ReportSyncMiss()
			if err != nil {
				return err
			}
		}
	} else {
		log.Infof("Content archival request %s is backing off, using cached video list", dlReq.Id)
	}

	//videos, err := dlReq.FetchVideoList()
	//if err != nil {
	//	return err
	//}
	//
	//log.Infof("Downloading %d videos for content type %s content value %s", len(videos), dlReq.ContentType, dlReq.ContentValue)

}
