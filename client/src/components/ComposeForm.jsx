import React, { useState } from 'react'
import { useSelector} from 'react-redux'
import api from "../api"

const ComposeForm = ({ callback }) => {
	const [form, setForm] = useState({
		content: "",
	})
	const [errors, setErrors] = useState([])

	const auth = useSelector(state => state.auth)

	const handleChange = (event) => {
        setForm({...form, [event.target.name]: event.target.value})
    }

	function handleSubmit(e) {
        e.preventDefault()

		api.post("tweets", form).then(() => {
			setForm({content: ""})
			setErrors([])

			callback()
		}).catch(error => {
			if ([400, 422].includes(error.response.status)) {
				setErrors(error.response.data.errors)
			}
		})
    }

	return (
		<div>
			{ errors.length > 0 &&
                <ul className="bg-red-100 px-3 py-2 my-4 rounded text-red-400 text-sm">
                    { errors.map((error) => <li key={error}>{error}</li>) }
                </ul>
            }

			<form onSubmit={handleSubmit} className="my-4">
				<textarea
					rows={3}
					name="content"
					className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
					placeholder={`What's on your mind, ${auth.user.username}?`}
					onChange={handleChange}
					value={form.content}
				/>
				<div className="mt-2 flex justify-end">
					<button
						type="submit"
						className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
					>
						Post
					</button>
				</div>
			</form>
		</div>
	)
}

export default ComposeForm