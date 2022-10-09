import { type IDownload } from "./types";
import LinearProgress from '@mui/material/LinearProgress';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { blockComponent } from "#components/meta";
import {
  Card,
  CardBody,
  CardHeader,
  DL,
  DS,
  type ICardProps,
} from "#components/lists";
import { Heading } from "#components/headings";

export interface IDownloadCardProps extends ICardProps {
  download: IDownload;
}

export const DownloadCard = blockComponent(undefined, Component);

function Component({
  download,
  headingLevel,
  ...blockProps
}: IDownloadCardProps) {
  const { progress, DlStatus, VideoID, Website } = download;

  return (
    <Card {...blockProps}>
      <CardHeader>
        <Heading level={headingLevel}>Video {VideoID}</Heading>
      </CardHeader>
      <CardBody>
        <DL>
          <DS isHorizontal dKey="Website" dValue={Website} />
          <DS isHorizontal dKey="Status" dValue={DlStatus} />
          <p></p>
         <span>Progress: <LinearProgress color="success" variant="determinate" value={progress} />{progress}%</span>

        </DL>
      </CardBody>
    </Card>
  );
}
