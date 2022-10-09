export default interface IDownload {
  VideoID: string;
  progress: number;
  DlStatus: "Downloading" | "Queued";
  Website: string;
}

export interface Deletion {
  Type: "deletion";
  Video: {
    VideoID: string;
    Website: string;
    DlStatus: "";
  };
}

export interface InProgress {
  total_bytes: number;
  downloaded_bytes: number;
  info_dict: {
    id: string;
  };
}

export interface Insertion {
  Type: "insertion";
  Video: {
    VideoID: string;
    Website: string;
    DlStatus: "Queued" | "Downloading";
  };
}
