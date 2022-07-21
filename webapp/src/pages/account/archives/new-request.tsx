import { useState } from "react";
import { Table, Button, Space } from "antd";

import styles from "./new-request.module.scss";

import { Page } from "#components/page";
import { NewVideoForm } from "#components/posts";
import { IArchivalRequest } from "#codegen/schema/001_interfaces";
import { ArchiveStatus } from "#components/archives";
import { deleteArchivalRequest, retryArchivalRequest } from "#api/archives";

function NewArchivePage() {
  const [archivalSubscriptions, setArchivalSubscriptions] = useState<
    IArchivalRequest[]
  >([]);

  const columns = [
    {
      title: "Status",
      key: "Status",
      render: (
        text: string,
        record: {
          ArchivedVideos: number;
          CurrentTotalVideos: number;
          LastSynced: undefined;
        }
      ) => (
        <span>
          <ArchiveStatus record={record} />
        </span>
      ),
    },
    {
      title: "URL",
      dataIndex: "Url",
    },
    {
      title: "Last synced",
      dataIndex: "LastSynced",
    },
    // {
    //     title: 'Days until next sync',
    //     'dataIndex': 'BackoffFactor',
    //     key: 'BackoffFactor',
    // },
    {
      title: "Downloaded",
      key: "Downloaded",
      render: (
        text: string,
        record: {
          ArchivedVideos: number;
          CurrentTotalVideos: number;
          LastSynced: null;
          UndownloadableVideos: number;
        }
      ) => (
        <span>
          <b>{record.ArchivedVideos + "/" + record.CurrentTotalVideos}</b>{" "}
          videos ({record.UndownloadableVideos} undownloadable)
        </span>
      ),
    },
    {
      title: "Actions",
      key: "Actions",
      render: (
        text: string,
        record: {
          ArchivedVideos: number;
          CurrentTotalVideos: number;
          LastSynced: null;
          DownloadID: number;
        }
      ) => (
        <Space size="middle">
          <Button
            className="background-blue"
            onClick={() => retryRequest(record.DownloadID)}
          >
            Retry {record.DownloadID}
          </Button>
          <Button
            className="background-blue"
            onClick={() => deleteRequest(record.DownloadID)}
          >
            Delete {record.DownloadID}
          </Button>
        </Space>
      ),
    },
  ];

  async function deleteRequest(download_id: number) {
    const formParams = new URLSearchParams([
      ["download_id", String(download_id)],
    ]);
    await deleteArchivalRequest(formParams);
  }

  async function retryRequest(download_id: number) {
    const formParams = new URLSearchParams([
      ["download_id", String(download_id)],
    ]);
    await retryArchivalRequest(formParams);
  }

  async function createNewArchival(url: string) {
    const subs = archivalSubscriptions ? archivalSubscriptions : [];
    const newList = [
      { Url: url, ArchivedVideos: 0, CurrentTotalVideos: 0, BackoffFactor: 1 },
      ...subs,
    ];
    setArchivalSubscriptions(newList);
  }

  return (
    <Page>
      <h1 className={styles.heading}>View and manage your archives</h1>
      <div className={styles.table}>
        <NewVideoForm onNewURL={createNewArchival} />
        <Table
          // @ts-expect-error types
          dataSource={archivalSubscriptions}
          scroll={{ y: 700 }}
          className="align-bottom w-full "
          ellipsis={true}
          // @ts-expect-error types
          columns={columns}
        />
      </div>
    </Page>
  );
}

export default NewArchivePage;
