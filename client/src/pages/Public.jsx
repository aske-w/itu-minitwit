import React, { useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import Tweets from '../components/Tweets';
import api from '../api';
import ComposeForm from '../components/ComposeForm';

const Public = () => {
    const [tweets, setTweets] = useState([])
    const auth = useSelector(state => state.auth)

    const fetchTweets = () => {
        api.get("/tweets")
            .then(response => {
                setTweets(response.data)
            })
    }

    useEffect(() => {
        fetchTweets()
    },[])

    return (
        <div>
            <h2 className="text-2xl font-semibold mb-4">Public timeline</h2>

            { auth.isLoggedIn && <ComposeForm callback={fetchTweets}/> }

            <Tweets tweets={tweets}/>
        </div>
    )
}

export default Public