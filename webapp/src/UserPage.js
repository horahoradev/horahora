import { useCallback, useEffect, useRef, useState } from "react";
import { Tag, Avatar, Button } from "antd";
import { Link, useParams } from "react-router-dom";
import dashjs from "dashjs";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUserCircle } from "@fortawesome/free-solid-svg-icons";
import Footer from "./Footer";
import { UserOutlined } from '@ant-design/icons';

import * as API from "./api";
import { Header } from "./components/header";
import { UserRank } from "./api/types";
import VideoList from "./VideoList";
import Paginatione from "./Pagination";

// {"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":11,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}

function UserPage() {
    let { id } = useParams();

    const [userData, setUserData] = useState(null);
    const [pageUserData, setPageUserData] = useState({});
    const [currPage, setPage] = useState(1);

    // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
    useEffect(() => {
        let ignore = false;

        let fetchData = async () => {
            let pageUserData = await API.getUser(id, currPage);
            if (!ignore) setPageUserData(pageUserData);
            
            let userData = await API.getUserdata();
            if (!ignore) setUserData(userData);
        };

        fetchData();
        return () => {
            ignore = true;
        };
    }, [id, currPage]);


    return (
        <>
            <Header userData={userData} />
            <div className="pt-10">
                {/*lol oh no*/}
                {/* TODO: add user profile image*/}
                <div className={"flex justify-center mx-4"}>
                    <Avatar shape="square" size={96} icon={<UserOutlined />} />
                </div>

                <h1 className={"flex justify-center text-black"}>{pageUserData.Username} <b> {pageUserData.banned && "(banned)"} </b>
                </h1>
                {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && <h1 className={"flex justify-center"}><Button type="primary" onClick={() => API.banUser(pageUserData.UserID)}>Ban</Button></h1>}
                {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && <h1 className={"flex justify-center"}><Button type="primary" onClick={() => API.setUserMod(pageUserData.UserID)}>Promote to mod</Button></h1>}
                {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && <h1 className={"flex justify-center"}><Button type="primary" onClick={() => API.setUserAdmin(pageUserData.UserID)}>Promote to admin</Button></h1>}

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