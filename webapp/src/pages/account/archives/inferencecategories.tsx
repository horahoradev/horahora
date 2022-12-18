import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { Page } from "#components/page";
import { getArchivalEvents, getArchivalRequests, approveVideo, getUnapprovedVideos, unapproveVideo, getInferenceCategories } from "#api/archives";
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
import { NewInferenceForm } from "#components/posts/new-inference-category";

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

function InferenceCategoriesPage(params: any) {
  const [inferenceCategories, changeInferenceCategories] = useState([]);
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);

  const router = useRouter();
  const { query, isReady } = router;


  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  async function createNewInference(url: string) {
    reloadPage();
  }


  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    if (!isReady) {
      return
    }

    let fetchData = async () => {
      let resp = await getInferenceCategories();
        console.log(resp);
      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore && resp.Entries != undefined) {
        changeInferenceCategories(resp.Entries);
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
    <NewInferenceForm onNewCategory={createNewInference}></NewInferenceForm>
    <CssBaseline />
    <TableContainer component={Paper}>
    <Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
        <TableHead>
        <TableRow>
            <TableCell align="left">Tag</TableCell>
            <TableCell align="left">Category</TableCell>
        </TableRow>
        </TableHead>
        <TableBody>
        {inferenceCategories.length != 0 ? inferenceCategories.map((row: any) => (
            <TableRow
            key={row.VideoID}
            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
            <TableCell align="left">{row.Tag}</TableCell>
            <TableCell align="left">{row.Category}</TableCell>
            </TableRow>
        )): null }
        </TableBody>
    </Table>
    </TableContainer>
    </ThemeProvider>
  )
}

export default InferenceCategoriesPage;
