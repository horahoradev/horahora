import { useCallback, useEffect, useRef, useState } from "react";
import { Tag, Avatar, Button } from "antd";
import { Link, useParams } from "react-router-dom";
import dashjs from "dashjs";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUserCircle } from "@fortawesome/free-solid-svg-icons";

import * as API from "./api";
import Header from "./Header";
import { UserRank } from "./api/types";
import VideoList from "./VideoList";
import Paginatione from "./Pagination";

// {"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":11,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}


function ArchivalPage() {
    const [userData, setUserData] = useState(null);
    const [archivalSubscriptions, setArchivalSubscriptions] = useState([]);

    function createNewArchival() {
        const url = document.getElementById('url').value;
        API.postArchival(url);
        let newList = [...archivalSubscriptions, url];
        setArchivalSubscriptions(newList);
    }

    // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
    useEffect(() => {
        let ignore = false;

        let fetchData = async () => {
            let userData = await API.getUserdata();
            if (!ignore) setUserData(userData);

            let subscriptionData = await API.getArchivalSubscriptions();
            if (!ignore) setArchivalSubscriptions(subscriptionData.ArchivalRequests);
        };

        fetchData();
        return () => {
            ignore = true;
        };
    }, []);

    let elements = [];
    if (archivalSubscriptions) {
        elements = [
            archivalSubscriptions.map((subscription, idx) => <p>{subscription.url}</p>),
        ];
    }

    return (
        <>
            <Header userData={userData} />
            <br></br>
            Type URL to archive here: <input type="text" id="url">
            </input>
            <p></p>
            <Button className="background-blue" onClick={createNewArchival}>Submit</Button>
            <p></p>
            {elements}
            </>
    );
}

export default ArchivalPage;