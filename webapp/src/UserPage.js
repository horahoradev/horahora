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

function UserPage() {
    let { id } = useParams();

    const [pageData, setPageData] = useState(null);
    const [userData, setUserData] = useState(null);
    const [pageUserData, setPageUserData] = useState({});
    const [currPage, setPage] = useState(1);

    // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
    useEffect(() => {
        let ignore = false;

        let fetchData = async () => {
            let data = await API.getVideo(id);
            if (!ignore) setPageData(data);

            let userData = await API.getUserdata();
            if (!ignore) setUserData(userData);

            let pageUserData = await API.getUser(id, currPage);
            if (!ignore) setPageUserData(pageUserData);
        };

        fetchData();
        return () => {
            ignore = true;
        };
    }, [id, currPage]);

    if (pageData == null) return null;

    return (
        <>
            <Header userData={userData} />
            <div>
                {/*lol oh no*/}
                {/* TODO: add user profile image*/}
                <div className={"flex justify-center mx-4"}>
                    <img src="http://localhost:8082/static/images/placeholder.png" className={"rounded-full"}>
                    </img>
                    <h1>{pageUserData.Username}</h1>
                </div>
                {/* TODO: no more copy pasta! */}
                <div className="flex justify-center mx-4 min-h-screen">
                    <div className="max-w-screen-lg w-screen my-6">
                        <VideoList videos={pageUserData.Videos} />
                        <Paginatione paginationData={pageUserData.PaginationData} onPageChange={setPage}/>
                    </div>
                </div>
            </div>
        </>
    );
}

export default UserPage;