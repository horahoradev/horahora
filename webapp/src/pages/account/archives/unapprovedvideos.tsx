import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { Page } from "#components/page";
import { getArchivalEvents, getArchivalRequests, approveVideo, getUnapprovedVideos, unapproveVideo } from "#api/archives";
import { type IArchivalEvent } from "#codegen/schema/001_interfaces";
import { CardList } from "#components/lists";
import { LoadingBar } from "#components/loading-bar";
import { EventCard } from "#entities/event";
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Button } from "#components/buttons";

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

function UnapprovedVideosPage(params: any) {
  const [unapprovedvideos, changeUnapprovedVideos] = useState<[]>();
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);

  const router = useRouter();
  const { query, isReady } = router;


  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  async function approveNewVideo(videoID: string) {
        await approveVideo(videoID);

        // quite lazy, forces a reload
        setTimerVal((timerVal) => timerVal + 1);
    }

  async function unapproveNewVideo(videoID: string) {
      await unapproveVideo(videoID);

      // quite lazy, forces a reload
      setTimerVal((timerVal) => timerVal + 1);
  }


  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    if (!isReady) {
      return
    }

    let fetchData = async () => {
      let resp = await getUnapprovedVideos();

      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore) {
        changeUnapprovedVideos(resp);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [timerVal, isReady]);

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);


  return (
    <ThemeProvider theme={darkTheme}>
    <CssBaseline />
    <TableContainer component={Paper}>
    <Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
        <TableHead>
        <TableRow>
            <TableCell align="left">ID</TableCell>
            <TableCell align="left">URL</TableCell>
            <TableCell align="left">Approve</TableCell>
            <TableCell align="left">Unapprove</TableCell>
        </TableRow>
        </TableHead>
        <TableBody>
        {unapprovedvideos ? unapprovedvideos.map((row: any) => (
            <TableRow
            key={row.VideoID}
            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
            <TableCell align="left">{row.VideoID}</TableCell>
            <TableCell align="left">{<a href={row.URL}>{row.URL}</a>}</TableCell>
            <TableCell align="left"><Button onClick={()=>approveNewVideo(row.VideoID)}></Button></TableCell>
            <TableCell align="left"><Button onClick={()=>unapproveNewVideo(row.VideoID)}></Button></TableCell>

            </TableRow>
        )): null }
        </TableBody>
    </Table>
    </TableContainer>
    </ThemeProvider>
  )
}

export default UnapprovedVideosPage;
