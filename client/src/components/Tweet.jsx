import React from 'react'

const Tweet = ({ tweet }) => {
	return (
		<li>
			<div className="flex space-x-3 border bg-white px-3 py-3 rounded-lg shadow-sm">
				<div>
					<img src={tweet.Gravatar_Url} alt="profile_img" />
				</div>

				<div className="flex-1">
					<div className="flex justify-between mb-1">
						<div className="text-sm font-semibold text-gray-800">{tweet.Username}</div>
						<div className="text-sm text-gray-400">{tweet.Format_Datetime}</div>
					</div>

					<div className="text-sm text-gray-600 line-clamp-2">{tweet.Text}</div>
				</div>
			</div>
		</li>
	)
}

export default Tweet