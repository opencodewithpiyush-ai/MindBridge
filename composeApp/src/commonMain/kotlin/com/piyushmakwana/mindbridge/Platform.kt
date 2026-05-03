package com.piyushmakwana.mindbridge

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform