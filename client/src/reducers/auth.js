import { createSlice } from '@reduxjs/toolkit'

export const authSlice = createSlice({
    name: "authStore",

    initialState:{
        user: JSON.parse(localStorage.getItem("user")),
        isLoggedIn: !! JSON.parse(localStorage.getItem("user")) || false
    },

    reducers: {
        login: (state, { payload }) => {
            state.user = payload
            state.isLoggedIn = true

            localStorage.setItem("user", JSON.stringify(payload))
        },

        logout: state => {
            state.user = null
            state.isLoggedIn = false

            localStorage.removeItem("user")
        },
    }
})

export const { login, logout } = authSlice.actions

export default authSlice.reducer