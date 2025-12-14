# Media

> Listen, um, I don’t really know what’s going on here, but I’m leaving. I
> don’t know where exactly, but I’m gonna start over. Come with me.
>
> [Paint it Black][1] (2016)

Go packages for streaming API authentication, metadata, and DRM playback.

**Key functionalities implemented across the packages include:**

- **Authentication & Session Management:** Methods to handle user logins,
   device linking (e.g., Roku, HBO Max), token exchanges, and session
   refreshing. It handles various authentication schemes including OAuth, JWT, and
   cookie-based sessions.

- **Content Discovery & Metadata:** Functions to resolve URLs (slugs) to
   internal content IDs, fetch series/episode details via JSON or GraphQL
   endpoints, and extract available video qualities.

- **Playback Extraction:** Logic to retrieve streaming manifests, specifically
   **MPEG-DASH (.mpd)** files, for movies and TV shows.

- **DRM Licensing:** implementations for interacting with Digital Rights
   Management systems (primarily **Widevine** and **PlayReady**). This includes
   generating challenge payloads, signing requests (HMAC/AES), and communicating
   with license servers to authorize playback.

- **HTTP Client Customization:** configuration of HTTP requests with specific
   headers (User-Agents, platform identifiers, custom tokens) required to
   masquerade as legitimate client devices (e.g., Android, Web, TV apps).

**Supported Services identified in the file structure include:**

1. PlutoTV: Live TV & Free Movies
2. Tubi: Free Movies & Live TV
3. HBO Max: Stream TV & Movies
4. Hulu: Stream TV shows & movies
5. Plex: Stream Movies & TV
6. CANAL+, Live and catch-up TV
7. Paramount+
8. ITVX
9. The NBC App - Stream TV Shows
10. Molotov - TV en direct, replay
11. MUBI: Curated Cinema
12. Rakuten TV -Movies & TV Series
13. CTV
14. Kanopy
15. AMC+
16. RTBF Auvio : direct et replay
17. The Roku Channel
18. The Criterion Channel
19. Draken Film
20. CineMember

[1]://f002.backblazeb2.com/file/minerals/Paint.It.Black.2016.mp4

## contact

<dl>
   <dt>email</dt>
      <dd>27@riseup.net</dd>
   <dt>Discord username</dt>
      <dd>10308</dd>
   <dt>PayPal</dt>
      <dd>https://paypal.com/donate?hosted_button_id=59UKABTT2F8LS</dd>
</dl>
