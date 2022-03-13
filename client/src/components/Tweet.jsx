import React from 'react'

const Tweet = ({ tweet }) => {
  return (
      <li>
        <img src={tweet.Gravatar_Url} alt="profile_img" />
        <p>
            <strong>
                <a href="/{Username}">{tweet.Username}</a>
            </strong>
            <span>{tweet.Text}</span>
            <small>&mdash;{tweet.Format_Datetime}</small>
        </p>
    </li>
  )
}

export default Tweet