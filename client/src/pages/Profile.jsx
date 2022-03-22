import React, { useEffect, useState } from 'react'
import Tweets from '../components/Tweets';
import { useSelector } from 'react-redux'
import { useNavigate, useParams } from "react-router-dom";
import api from '../api';

const Profile = (props) => {
    const [user, setUser] = useState(null)
    const [tweets, setTweets] = useState([])
    const [isFollowing, setIsFollowing] = useState(false)
    const auth = useSelector(state => state.auth)
    const navigate = useNavigate()
    const { username } = useParams()

    const fetchUser = () => {
        api.get(`/users/${username}`)
            .then(response => {
                setUser(response.data)
            })
            .catch(error => {
                navigate("/public")
            })
        // Checks if user is already following or not
        api.get(`/users/${username}/isfollowing`).then(response => {
                setIsFollowing(response.data.isFollowing)
            })
    }

    const fetchTweets = () => {
        api.get(`/users/${user.username.toLowerCase()}/tweets`)
            .then(response => {
                setTweets(response.data)
            })
    }

    const handleFollow = (e) => {
        e.preventDefault()
        api.post(`/users/${username}/follow`, { username }).then(() => {
            setIsFollowing(prevValue => !prevValue)
        })

    }

    useEffect(() => {
        fetchUser()
    }, [])

    useEffect(() => {
        if ( user) {
            fetchTweets()
        }
    }, [user])



    if ( ! user) {
        return null
    }

    return (
        <div>
            <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-semibold">{ user.username }'s timeline</h2>

                { auth.isLoggedIn &&
                    <button
                        type="button"
                        className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                        onClick={handleFollow}
                    >
                        {isFollowing ? "Unfollow" : "Follow"}
                    </button>
                }
            </div>

            <Tweets tweets={tweets}/>
        </div>
    )
}

export default Profile