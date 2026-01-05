package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golift.io/starr"
	"golift.io/starr/sonarr"

	"github.com/kyleaupton/snaggle/backend/internal/release"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// testTitles contains release titles to test against Sonarr
// These are used to validate parity between Snaggle and Sonarr's quality parser
var testTitles = []string{
	// SDTV
	"S07E23 .avi",
	"The.Series.S01E13.x264-CtrlSD",
	"The Series S02E01 HDTV XviD 2HD",
	"The Series S05E11 PROPER HDTV XviD 2HD",
	"The Series Show S02E08 HDTV x264 FTP",
	"The.Series.2011.S02E01.WS.PDTV.x264-TLA",
	"The.Series.2011.S02E01.WS.PDTV.x264-REPACK-TLA",
	"The Series S01E04 DSR x264 2HD",
	"The Series S01E04 Series Death Train DSR x264 MiNDTHEGAP",
	"The Series S11E03 has no periods or extension HDTV",
	"The.Series.S04E05.HDTV.XviD-LOL",
	"The.Series.S02E15.avi",
	"The.Series.S02E15.xvid",
	"The.Series.S02E15.divx",
	"The.Series.S03E06.HDTV-WiDE",
	"Series.S10E27.WS.DSR.XviD-2HD",
	"[HorribleSubs] The Series - 32 [480p]",
	"[CR] The Series - 004 [480p][48CE2D0F]",
	"[Hatsuyuki] The Series - 363 [848x480][ADE35E38]",
	"The.Series.S03.TVRip.XviD-NOGRP",
	"[HorribleSubs] The Series - 03 [360p].mkv",
	"[SubsPlease] Series Title (540p) [AB649D32].mkv",
	"[Erai-raws] Series Title [540p][Multiple Subtitle].mkv",

	// DVD
	"The.Series.S01E13.NTSC.x264-CtrlSD",
	"The.Series.S03E06.DVDRip.XviD-WiDE",
	"The.Series.S03E06.DVD.Rip.XviD-WiDE",
	"the.Series.1x13.circles.ws.xvidvd-tns",
	"the_Series.9x18.sunshine_days.ac3.ws_dvdrip_xvid-fov.avi",
	"[FroZen] Series - 23 [DVD][7F6170E6]",
	"[AniDL] Series - 26 -[360p][DVD][D - A][Exiled - Destiny]",

	// WEBDL-480p
	"The.Series.S01E10.The.Leviathan.480p.WEB-DL.x264-mSD",
	"The.Series.S04E10.Glee.Actually.480p.WEB-DL.x264-mSD",
	"The.SeriesS06E11.The.Santa.Simulation.480p.WEB-DL.x264-mSD",
	"The.Series.S02E04.480p.WEB.DL.nSD.x264-NhaNc3",
	"The.Series.S01E08.Das.geloeschte.Ich.German.Dubbed.DL.AmazonHD.x264-TVS",
	"The.Series.S01E04.Rod.Trip.mit.meinem.Onkel.German.DL.NetflixUHD.x264",
	"[HorribleSubs] Series Title! S01 [Web][MKV][h264][480p][AAC 2.0][Softsubs (HorribleSubs)]",
	"Series.Title.S13E11.Ausgebacken.German.AmazonSD.h264-4SF",

	// Bluray-480p
	"SERIES.S03E01-06.DUAL.XviD.Bluray.AC3-REPACK.-HELLYWOOD.avi",
	"SERIES.S03E01-06.DUAL.BDRip.XviD.AC3.-HELLYWOOD",
	"SERIES.S03E01-06.DUAL.BDRip.X-viD.AC3.-HELLYWOOD",
	"SERIES.S03E01-06.DUAL.BDRip.AC3.-HELLYWOOD",
	"SERIES.S03E01-06.DUAL.BDRip.XviD.AC3.-HELLYWOOD.avi",
	"SERIES.S03E01-06.DUAL.XviD.Bluray.AC3.-HELLYWOOD.avi",
	"The.Series.S01E05.480p.BluRay.DD5.1.x264-HiSD",
	"The Series (BD)(640x480(RAW) (BATCH 1) (1-13)",
	"[Doki] Series - 02 (848x480 XviD BD MP3) [95360783]",
	"Adventures.of.Sonic.the.Hedgehog.S01.BluRay.480i.DD.2.0.AVC.REMUX-FraMeSToR",
	"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.480i.DD.2.0.AVC.REMUX-FraMeSToR",

	// WEBRip-480p
	"The.Series.S02E10.480p.HULU.WEBRip.x264-Puffin",
	"The.Series.S10E14.Techs.And.Balances.480p.AE.WEBRip.AAC2.0.x264-SEA",
	"Series.Title.1x04.ITA.WEBMux.x264-NovaRip",

	// Bluray-576p
	"The.Series.S01E05.576p.BluRay.DD5.1.x264-HiSD",

	// HDTV-720p
	"Series - S01E01 - Title [HDTV]",
	"Series - S01E01 - Title [HDTV-720p]",
	"The Series S04E87 REPACK 720p HDTV x264 aAF",
	"The.Series.S02E15.720p",
	"S07E23 - [HDTV-720p].mkv",
	"Series - S22E03 - MoneyBART - HD TV.mkv",
	"S07E23.mkv",
	"The.Series.S08E05.720p.HDTV.X264-DIMENSION",
	"The.Series.S02E15.mkv",
	"The.Series.S01E08.Tourmaline.Nepal.720p.HDTV.x264-DHD",
	"[Underwater-FFF] The Series - 01 (720p) [27AAA0A0]",
	"[Doki] The Series - 07 (1280x720 Hi10P AAC) [80AF7DDE]",
	"[Doremi].The.Series.5.Go.Go!.31.[1280x720].[C65D4B1F].mkv",
	"[HorribleSubs]_Series_Title_-_145_[720p]",
	"[Eveyuu] Series Title - 10 [Hi10P 1280x720 H264][10B23BD8]",
	"The.Series.US.S12E17.HR.WS.PDTV.X264-DIMENSION",
	"The.Series.The.Lost.Sonarr.Summer.HR.WS.PDTV.x264-DHD",
	"The Series S01E07 - Motor zmen (CZ)[TvRip][HEVC][720p]",
	"The.Series.S05E06.720p.HDTV.x264-FHD",
	"Series.Title.1x01.ITA.720p.x264-RlsGrp [01/54] - \"series.title.1x01.ita.720p.x264-rlsgrp.nfo\"",
	"[TMS-Remux].Series.Title.X.21.720p.[76EA1C53].mkv",

	// HDTV-1080p
	"Under the Series S01E10 Let the Sonarr Begin 1080p",
	"Series.S07E01.ARE.YOU.1080P.HDTV.X264-QCF",
	"Series.S07E01.ARE.YOU.1080P.HDTV.x264-QCF",
	"Series.S07E01.ARE.YOU.1080P.HDTV.proper.X264-QCF",
	"Series - S01E01 - Title [HDTV-1080p]",
	"[HorribleSubs] Series Title - 32 [1080p]",
	"Series S01E07 - Sonarr zmen (CZ)[TvRip][HEVC][1080p]",
	"The Online Series Alicization 04 vostfr FHD",
	"Series Slayer 04 vostfr FHD.mkv",
	"[Onii-ChanSub] The.Series - 02 vostfr (FHD 1080p 10bits).mkv",
	"[Miaou] Series Title 02 VOSTFR FHD 10 bits",
	"[mhastream.com]_Episode_05_FHD.mp4",
	"[Kousei]_One_Series_ - _609_[FHD][648A87C7].mp4",
	"Series culpable 1x02 Culpabilidad [HDTV 1080i AVC MP2 2.0 Sub][GrupoHDS]",
	"Series como paso - 19x15 [344] Cuarenta anos de baile [HDTV 1080i AVC MP2 2.0 Sub][GrupoHDS]",
	"Super.Seires.Go.S01E02.Depths.of.Sonarr.1080i.HDTV.DD5.1.H.264-NOGRP",

	// HDTV-2160p
	"My Title - S01E01 - EpTitle [HEVC 4k DTSHD-MA-6ch]",
	"My Title - S01E01 - EpTitle [HEVC-4k DTSHD-MA-6ch]",
	"My Title - S01E01 - EpTitle [4k HEVC DTSHD-MA-6ch]",

	// WEBDL-720p
	"Series S01E04 Mexicos Death Train 720p WEB DL",
	"Series Five 0 S02E21 720p WEB DL DD5 1 H 264",
	"Series S04E22 720p WEB DL DD5 1 H 264 NFHD",
	"Series - S11E06 - D-Yikes! - 720p WEB-DL.mkv",
	"The.Series.S02E15.720p.WEB-DL.DD5.1.H.264-SURFER",
	"S07E23 - [WEBDL].mkv",
	"Series S04E22 720p WEB-DL DD5.1 H264-EbP.mkv",
	"Series.S04.720p.Web-Dl.Dd5.1.h264-P2PACK",
	"Da.Series.Shows.S02E04.720p.WEB.DL.nSD.x264-NhaNc3",
	"Series.Miami.S04E25.720p.iTunesHD.AVC-TVS",
	"Series.S06E23.720p.WebHD.h264-euHD",
	"Series.Title.2016.03.14.720p.WEB.x264-spamTV",
	"Series.Title.2016.03.14.720p.WEB.h264-spamTV",
	"Series.S01E08.Das.geloeschte.Ich.German.DD51.Dubbed.DL.720p.AmazonHD.x264-TVS",
	"Series.Polo.S01E11.One.Hundred.Sonarrs.2015.German.DD51.DL.720p.NetflixUHD.x264.NewUp.by.Wunschtante",
	"Series 2016 German DD51 DL 720p NetflixHD x264-TVS",
	"Series.6x10.Basic.Sonarr.Repair.and.Replace.ITA.ENG.720p.WEB-DLMux.H.264-GiuseppeTnT",
	"Series.6x11.Modern.Spy.ITA.ENG.720p.WEB.DLMux.H.264-GiuseppeTnT",
	"The Series Was Dead 2010 S09E13 [MKV / H.264 / AC3/AAC / WEB / Dual Audio / Ingles / 720p]",
	"into.the.Series.s03e16.h264.720p-web-handbrake.mkv",
	"Series.S01E01.The.Sonarr.Principle.720p.WEB-DL.DD5.1.H.264-BD",
	"Series.S03E05.Griebnitzsee.German.720p.MaxdomeHD.AVC-TVS",
	"[HorribleSubs] Series Title! S01 [Web][MKV][h264][720p][AAC 2.0][Softsubs (HorribleSubs)]",
	"[HorribleSubs] Series Title! S01 [Web][MKV][h264][AAC 2.0][Softsubs (HorribleSubs)]",
	"Series.Title.S04E13.960p.WEB-DL.AAC2.0.H.264-squalor",
	"Series.Title.S16.DP.WEB.720p.DDP.5.1.H.264.PLEX",
	"Series.Title.S01E01.Erste.Begegnungen.German.DD51.Synced.DL.720p.HBOMaxHD.AVC-TVS",
	"Series.Title.S01E05.Tavora.greift.an.German.DL.720p.DisneyHD.h264-4SF",

	// WEBRip-720p
	"Series.Title.S04E01.720p.WEBRip.AAC2.0.x264-NFRiP",
	"Series.Title.S01E07.A.Prayer.For.Mad.Sweeney.720p.AMZN.WEBRip.DD5.1.x264-NTb",
	"Series.Title.S07E01.A.New.Home.720p.DSNY.WEBRip.AAC2.0.x264-TVSmash",
	"Series.Title.1x04.ITA.720p.WEBMux.x264-NovaRip",

	// WEBDL-1080p
	"Series S09E03 1080p WEB DL DD5 1 H264 NFHD",
	"Two and a Half Developers of the Series S10E03 1080p WEB DL DD5 1 H 264 NFHD",
	"Series.S08E01.1080p.WEB-DL.DD5.1.H264-NFHD",
	"Its.Always.Sonarrs.Fault.S08E01.1080p.WEB-DL.proper.AAC2.0.H.264",
	"This is an Easter Egg S10E03 1080p WEB DL DD5 1 H 264 REPACK NFHD",
	"Series.S04E09.Swan.Song.1080p.WEB-DL.DD5.1.H.264-ECI",
	"The.Big.Easter.Theory.S06E11.The.Sonarr.Simulation.1080p.WEB-DL.DD5.1.H.264",
	"Sonarr's.Baby.S01E02.Night.2.[WEBDL-1080p].mkv",
	"Series.Title.2016.03.14.1080p.WEB.x264-spamTV",
	"Series.Title.2016.03.14.1080p.WEB.h264-spamTV",
	"Series.S01.1080p.WEB-DL.AAC2.0.AVC-TrollHD",
	"Series Title S06E08 1080p WEB h264-EXCLUSIVE",
	"Series Title S06E08 No One PROPER 1080p WEB DD5 1 H 264-EXCLUSIVE",
	"Series Title S06E08 No One PROPER 1080p WEB H 264-EXCLUSIVE",
	"The.Series.S25E21.Pay.No1.1080p.WEB-DL.DD5.1.H.264-NTb",
	"Series.S01E08.Das.geloeschte.Ich.German.DD51.Dubbed.DL.1080p.AmazonHD.x264-TVS",
	"Death.Series.2017.German.DD51.DL.1080p.NetflixHD.x264-TVS",
	"Series.S01E08.Pro.Gamer.1440p.BKPL.WEB-DL.H.264-LiGHT",
	"Series.Title.S04E11.Teddy's.Choice.FHD.1080p.Web-DL",
	"Series.S04E03.The.False.Bride.1080p.NF.WEB.DDP5.1.x264-NTb[rartv]",
	"Series.Title.S02E02.This.Year.Will.Be.Different.1080p.AMZN.WEB...",
	"Series.Title.S02E02.This.Year.Will.Be.Different.1080p.AMZN.WEB.",
	"Series Title - S01E11 2020 1080p Viva MKV WEB",
	"[HorribleSubs] Series Title! S01 [Web][MKV][h264][1080p][AAC 2.0][Softsubs (HorribleSubs)]",
	"[LostYears] Series Title - 01-17 (WEB 1080p x264 10-bit AAC) [Dual-Audio]",
	"Series.and.Titles.S01.1080p.NF.WEB.DD2.0.x264-SNEAkY",
	"Series.Title.S02E02.This.Year.Will.Be.Different.1080p.WEB.H 265",
	"Series Title Season 2 [WEB 1080p HEVC Opus] [Netaro]",
	"Series Title Season 2 (WEB 1080p HEVC Opus) [Netaro]",
	"Series.Title.S01E01.Erste.Begegnungen.German.DD51.Synced.DL.1080p.HBOMaxHD.AVC-TVS",
	"Series.Title.S01E05.Tavora.greift.an.German.DL.1080p.DisneyHD.h264-4SF",
	"Series.Title.S02E04.German.Dubbed.DL.AAC.1080p.WEB.AVC-GROUP",

	// WEBRip-1080p
	"Series.Title.S04E01.iNTERNAL.1080p.WEBRip.x264-QRUS",
	"Series.Title.S07E20.1080p.AMZN.WEBRip.DDP5.1.x264-ViSUM ac3.(NLsub)",
	"Series.Title.S03E09.1080p.NF.WEBRip.DD5.1.x264-ViSUM",
	"The Series 42 S09E13 1.54 GB WEB-RIP 1080p Dual-Audio 2019 MKV",
	"Series.Title.1x04.ITA.1080p.WEBMux.x264-NovaRip",
	"Series.Title.2019.S02E07.Chapter.15.The.Believer.4Kto1080p.DSNYP.Webrip.x265.10bit.EAC3.5.1.Atmos.GokiTAoE",
	"Series.Title.S01.1080p.AMZN.WEB-Rip.DDP5.1.H.264-Telly",

	// WEBDL-2160p
	"Series.Title.2016.03.14.2160p.WEB.x264-spamTV",
	"Series.Title.2016.03.14.2160p.WEB.h264-spamTV",
	"Series.Title.2016.03.14.2160p.WEB.PROPER.h264-spamTV",
	"House.of.Sonarr.AK.s05e13.4K.UHD.WEB.DL",
	"House.of.Sonarr.AK.s05e13.UHD.4K.WEB.DL",
	"[HorribleSubs] Series Title! S01 [Web][MKV][h264][2160p][AAC 2.0][Softsubs (HorribleSubs)]",
	"Series Title S02 2013 WEB-DL 4k H265 AAC 2Audio-HDSWEB",
	"Series.Title.S02E02.This.Year.Will.Be.Different.2160p.WEB.H.265",
	"Series.Title.S02E04.German.Dubbed.DL.AAC.2160p.DV.HDR.WEB.HEVC-GROUP",

	// WEBRip-2160p
	"Series S01E01.2160P AMZN WEBRIP DD2.0 HI10P X264-TROLLUHD",
	"JUST ADD SONARR S01E01.2160P AMZN WEBRIP DD2.0 X264-TROLLUHD",
	"The.Man.In.The.Series.S01E01.2160p.AMZN.WEBRip.DD2.0.Hi10p.X264-TrollUHD",
	"The Man In the Series S01E01 2160p AMZN WEBRip DD2.0 Hi10P x264-TrollUHD",
	"House.of.Sonarr.AK.S05E08.Chapter.60.2160p.NF.WEBRip.DD5.1.x264-NTb.NLsubs",
	"Sonarr Saves the World S01 2160p Netflix WEBRip DD5.1 x264-TrollUHD",

	// Bluray-720p
	"SERIES.S03E01-06.DUAL.Bluray.AC3.-HELLYWOOD.avi",
	"Series - S01E03 - Come Fly With Me - 720p BluRay.mkv",
	"The Big Series.S03E01.The Sonarr Can Opener.m2ts",
	"Series.S01E02.Chained.Sonarr.[Bluray720p].mkv",
	"[FFF] DATE A Sonarr Dev - 01 [BD][720p-AAC][0601BED4]",
	"[coldhell] Series v3 [BD720p][03192D4C]",
	"[RandomRemux] Series - 01 [720p BD][043EA407].mkv",
	"[Kaylith] Series Friends Specials - 01 [BD 720p AAC][B7EEE164].mkv",
	"SERIES.S03E01-06.DUAL.Blu-ray.AC3.-HELLYWOOD.avi",
	"SERIES.S03E01-06.DUAL.720p.Blu-ray.AC3.-HELLYWOOD.avi",
	"[Elysium]Lucky.Series.01(BD.720p.AAC.DA)[0BB96AD8].mkv",
	"Series.Galaxy.S01E01.33.720p.HDDVD.x264-SiNNERS.mkv",
	"The.Series.S01E07.RERIP.720p.BluRay.x264-DEMAND",
	"Sans.Series.De.Traces.FRENCH.720p.BluRay.x264-FHD",
	"Series.Black.1x01.Selezione.Naturale.ITA.720p.BDMux.x264-NovaRip",
	"Series.Hunter.S02.720p.Blu-ray.Remux.AVC.FLAC.2.0-SiCFoI",
	"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.720p.DD.2.0.AVC.REMUX-FraMeSToR",

	// Bluray-1080p
	"Series - S01E03 - Come Fly With Me - 1080p BluRay.mkv",
	"Sonarr.Of.Series.S02E13.1080p.BluRay.x264-AVCDVD",
	"Series.S01E02.Chained.Heat.[Bluray1080p].mkv",
	"[FFF] Series no Muromi-san - 10 [BD][1080p-FLAC][0C4091AF]",
	"[coldhell] Series v2 [BD1080p][5A45EABE].mkv",
	"[Kaylith] Series Friends Specials - 01 [BD 1080p FLAC][429FD8C7].mkv",
	"[Zurako] Log Series - 01 - The Sonarr (BD 1080p AAC) [7AE12174].mkv",
	"SERIES.S03E01-06.DUAL.1080p.Blu-ray.AC3.-HELLYWOOD.avi",
	"[Coalgirls]_Series!!_01_(1920x1080_Blu-ray_FLAC)_[8370CB8F].mkv",
	"Planet.Series.S01E11.Code.Deep.1080p.HD-DVD.DD.VC1-TRB",
	"Series Away(2001) Bluray FHD Hi10P.mkv",
	"S for Series 2005 1080p UHD BluRay DD+7.1 x264-LoRD.mkv",
	"Series.Title.2011.1080p.UHD.BluRay.DD5.1.HDR.x265-CtrlHD.mkv",
	"Fall.Of.The.Release.Groups.S02E13.1080p.BDLight.x265-AVCDVD",

	// Bluray-1080p Remux
	"Series!!! on ICE - S01E12[JP BD Remux][ENG subs]",
	"Series.Title.S01E08.The.Well.BluRay.1080p.AVC.DTS-HD.MA.5.1.REMUX-FraMeSToR",
	"Series.Title.2x11.Nato.Per.La.Truffa.Bluray.Remux.AVC.1080p.AC3.ITA",
	"Series.Title.2x11.Nato.Per.La.Truffa.Bluray.Remux.AVC.AC3.ITA",
	"Series.Title.S03E01.The.Calm.1080p.DTS-HD.MA.5.1.AVC.REMUX-FraMeSToR",
	"Series Title Season 2 (BDRemux 1080p HEVC FLAC) [Netaro]",
	"[Vodes] Series Title - Other Title (2020) [BDRemux 1080p HEVC Dual-Audio]",
	"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.1080p.DD.2.0.AVC.REMUX-FraMeSToR",
	"Series Title S01 2018 1080p BluRay Hybrid-REMUX AVC TRUEHD 5.1 Dual Audio-ZR-",
	"Series.Title.S01.2018.1080p.BluRay.Hybrid-REMUX.AVC.TRUEHD.5.1.Dual.Audio-ZR-",

	// Bluray-2160p
	"Series.Title.US.s05e13.4K.UHD.Bluray",
	"Series.Title.US.s05e13.UHD.4K.Bluray",
	"[DameDesuYo] Series Bundle - Part 1 (BD 4K 8bit FLAC)",
	"Series.Title.2014.2160p.UHD.BluRay.X265-IAMABLE.mkv",
	"Series.Title.S05EO1.Episode.Title.2160p.BDRip.AAC.7.1.HDR10.x265.10bit-Markll",

	// Bluray-2160p Remux
	"Series!!! on ICE - S01E12[JP BD 2160p Remux][ENG subs]",
	"Series.Title.S01E08.The.Sonarr.BluRay.2160p.AVC.DTS-HD.MA.5.1.REMUX-FraMeSToR",
	"Series.Title.2x11.Nato.Per.The.Sonarr.Bluray.Remux.AVC.2160p.AC3.ITA",
	"[Dolby Vision] Sonarr.of.Series.S07.MULTi.UHD.BLURAY.REMUX.DV-NoTag",
	"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.2160p.DD.2.0.AVC.REMUX-FraMeSToR",
	"Series Title S01 2018 2160p BluRay Hybrid-REMUX AVC TRUEHD 5.1 Dual Audio-ZR-",
	"Series.Title.S01.2018.2160p.BluRay.Hybrid-REMUX.AVC.TRUEHD.5.1.Dual.Audio-ZR-",

	// Raw-HD
	"POI S02E11 1080i HDTV DD5.1 MPEG2-TrollHD",
	"How I Met Your Developer S01E18 Nothing Good Happens After Sonarr 720p HDTV DD5.1 MPEG2-TrollHD",
	"The Series S01E11 The Finals 1080i HDTV DD5.1 MPEG2-TrollHD",
	"Series.Title.S07E11.1080i.HDTV.DD5.1.MPEG2-NTb.ts",
	"Game of Series S04E10 1080i HDTV MPEG2 DD5.1-CtrlHD.ts",
	"Series.Title.S02E05.1080i.HDTV.DD2.0.MPEG2-NTb.ts",
	"Show - S03E01 - Episode Title Raw-HD.ts",
	"Series.Title.S10E09.Title.1080i.UPSCALE.HDTV.DD5.1.MPEG2-zebra",
	"Series.Title.2011-08-04.1080i.HDTV.MPEG-2-CtrlHD",
}

func main() {
	// Get configuration from environment
	sonarrURL := os.Getenv("SONARR_URL")
	sonarrAPIKey := os.Getenv("SONARR_API_KEY")

	if sonarrURL == "" || sonarrAPIKey == "" {
		fmt.Printf("%sError: SONARR_URL and SONARR_API_KEY environment variables must be set%s\n", colorRed, colorReset)
		fmt.Println("\nUsage:")
		fmt.Println("  SONARR_URL=http://localhost:8989 SONARR_API_KEY=your-key go run ./cmd/quality-test")
		os.Exit(1)
	}

	// Initialize Sonarr client
	cfg := starr.New(sonarrAPIKey, sonarrURL, 60*time.Second)
	client := sonarr.New(cfg)

	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║           Quality Parser Comparison: Sonarr vs Snaggle           ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n\n", colorCyan, colorReset)

	var matches, mismatches, errors int
	var mismatchDetails []string

	ctx := context.Background()

	for i, title := range testTitles {
		fmt.Printf("[%3d/%3d] Testing: %s\n", i+1, len(testTitles), truncate(title, 60))

		// Get Sonarr's parsed quality
		sonarrQuality, err := getSonarrQuality(ctx, client, title)
		if err != nil {
			fmt.Printf("          %s⚠ Sonarr error: %v%s\n\n", colorYellow, err, colorReset)
			errors++
			continue
		}

		// Get Snaggle's parsed quality
		snaggleQuality := getSnaggleQuality(title)

		// Compare results
		if sonarrQuality == snaggleQuality {
			matches++
			fmt.Printf("          %s✓ MATCH: %s%s\n\n", colorGreen, sonarrQuality, colorReset)
		} else {
			mismatches++
			detail := fmt.Sprintf("Title: %s\n          Sonarr:  %s\n          Snaggle: %s", title, sonarrQuality, snaggleQuality)
			mismatchDetails = append(mismatchDetails, detail)
			fmt.Printf("          %s✗ MISMATCH%s\n", colorRed, colorReset)
			fmt.Printf("            Sonarr:  %s\n", sonarrQuality)
			fmt.Printf("            Snaggle: %s\n\n", snaggleQuality)
		}
	}

	// Print summary
	total := matches + mismatches
	parity := float64(matches) / float64(total) * 100

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════════%s\n", colorCyan, colorReset)
	fmt.Printf("%s                              SUMMARY                              %s\n", colorCyan, colorReset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════════%s\n\n", colorCyan, colorReset)

	fmt.Printf("Total tests:  %d\n", len(testTitles))
	fmt.Printf("Processed:    %d\n", total)
	if errors > 0 {
		fmt.Printf("%sErrors:       %d%s\n", colorYellow, errors, colorReset)
	}
	fmt.Printf("%sMatches:      %d (%.1f%%)%s\n", colorGreen, matches, parity, colorReset)
	fmt.Printf("%sMismatches:   %d (%.1f%%)%s\n", colorRed, mismatches, 100-parity, colorReset)

	if len(mismatchDetails) > 0 {
		fmt.Printf("\n%s── Mismatch Details ──────────────────────────────────────────────%s\n\n", colorYellow, colorReset)
		for i, detail := range mismatchDetails {
			fmt.Printf("%d. %s\n\n", i+1, detail)
		}
	}

	// Exit with error code if there were mismatches
	if mismatches > 0 {
		os.Exit(1)
	}
}

// getSonarrQuality calls Sonarr's parse API and extracts the quality name
func getSonarrQuality(ctx context.Context, client *sonarr.Sonarr, title string) (string, error) {
	input := &sonarr.ParseInput{Title: title}
	result, err := client.ParseContext(ctx, input)
	if err != nil {
		return "", err
	}

	if result == nil || result.ParsedEpisodeInfo == nil || result.ParsedEpisodeInfo.Quality == nil || result.ParsedEpisodeInfo.Quality.Quality == nil {
		return "Unknown", nil
	}

	return result.ParsedEpisodeInfo.Quality.Quality.Name, nil
}

// getSnaggleQuality uses Snaggle's quality parser
func getSnaggleQuality(title string) string {
	result := release.Parse(title)
	return result.Quality.Full()
}

// truncate shortens a string to maxLen characters
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
