import React, { useEffect, useState } from "react";
import { StompSubscription, type IMessage } from "@stomp/stompjs";
import { useMutex } from "react-context-mutex";

import { Page } from "#components/page";
import { NewVideoForm } from "#components/posts";
import { IArchivalRequest } from "#codegen/schema/001_interfaces";
import { getArchivalRequests } from "#api/archives";
import { StartCardList } from "#components/lists";
import { WSClient, WSConfig } from "#lib/fetch";
import { fetchDownloadsInProgress } from "#api/archives";
import { DownloadCard, IDownload} from "#entities/download";
import { LoadingBar } from "#components/loading-bar";
import { ArchivalDownloadsPage } from "#components/video";
import { RequestCard } from "#entities/request";
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Typography from '@mui/material/Typography';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import ArchivalEventsPage from "./events/[download_id]";

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});
interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
          <Typography>{children}</Typography>
      )}
    </div>
  );
}

function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

function NewArchivePage() {
  const [requests, changeRequests] = useState<IArchivalRequest[]>();
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);
  const [value, setValue] = React.useState(0);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let subscriptionData = await getArchivalRequests();

      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore) {
        changeRequests(subscriptionData.ArchivalRequests);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [timerVal]);

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  async function createNewArchival(url: string) {
    const subs = requests ? requests : [];
    reloadPage();
  }

  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };

  function a11yProps(index: number) {
    return {
      id: `simple-tab-${index}`,
      'aria-controls': `simple-tabpanel-${index}`,
    };
  }

  return (

    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Page title="View and manage your archives">
        <span><NewVideoForm onNewURL={createNewArchival} /></span>
        <Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
          <Tab label={<b>Archives</b>} {...a11yProps(0)} />
          <Tab label={<b>Video Downloads</b>} {...a11yProps(1)} />
          <Tab label={<b>Logs</b>} {...a11yProps(2)} />
        </Tabs>
        <TabPanel value={value} index={0}>
        <StartCardList>
          {!requests ? (
            <LoadingBar />
          ) : (
            requests.map((request) => (
              <RequestCard
                key={request.DownloadID}
                request={request}
                onUpdate={reloadPage}
              />
            ))
          )}
        </StartCardList>
        </TabPanel>
        <TabPanel value={value} index={1}>
        {<ArchivalDownloadsPage></ArchivalDownloadsPage>}
        </TabPanel>
        <TabPanel value={value} index={2}>
        <ArchivalEventsPage downloadID="all"></ArchivalEventsPage>
        </TabPanel>
      </Page>
    </ThemeProvider>
  );
}

export default NewArchivePage;
