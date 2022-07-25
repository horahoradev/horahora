import { Tag } from "antd";
import { CheckOutlined, SyncOutlined } from "@ant-design/icons";

import { type IArchivalRequest } from "#codegen/schema/001_interfaces";

export interface IArchiverRequestStatusProps {
  record: IArchivalRequest;
}

export function ArchiveRequestStatus({ record }: IArchiverRequestStatusProps) {
  if (
    record.ArchivedVideos == record.CurrentTotalVideos &&
    record.CurrentTotalVideos != 0
  ) {
    return (
      <Tag color="green" className="p-1 text-base" icon={<CheckOutlined />}>
        Complete
      </Tag>
    );
  } else if (record.CurrentTotalVideos == 0 || record.LastSynced == null) {
    return (
      <Tag color="blue" className="text-base" icon={<SyncOutlined spin />}>
        Fetching...
      </Tag>
    );
  } else {
    return (
      <Tag color="blue" className="text-base" icon={<SyncOutlined spin />}>
        {" "}
        Archiving...
      </Tag>
    );
  }
}
