package quality

import (
	"testing"
)

func TestQualityMethods(t *testing.T) {
	tests := []struct {
		name       string
		quality    Quality
		source     string
		resolution string
		isRemux    bool
	}{
		{"Bluray1080p", Bluray1080p, "BluRay", "1080p", false},
		{"Bluray1080pRemux", Bluray1080pRemux, "BluRay", "1080p", true},
		{"WEBDL2160p", WEBDL2160p, "WEB-DL", "2160p", false},
		{"HDTV720p", HDTV720p, "HDTV", "720p", false},
		{"SDTV", SDTV, "SDTV", "SD", false},
		{"DVD", DVD, "DVD", "SD", false},
		{"Unknown", Unknown, "Unknown", "Unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.quality.Source(); got != tt.source {
				t.Errorf("Quality.Source() = %v, want %v", got, tt.source)
			}
			if got := tt.quality.Resolution(); got != tt.resolution {
				t.Errorf("Quality.Resolution() = %v, want %v", got, tt.resolution)
			}
			if got := tt.quality.IsRemux(); got != tt.isRemux {
				t.Errorf("Quality.IsRemux() = %v, want %v", got, tt.isRemux)
			}
		})
	}
}

// TestParseQuality tests the quality parser against expected values
// These expected values are derived from Sonarr's QualityParser behavior
func TestParseQuality(t *testing.T) {
	tests := []struct {
		title    string
		expected Quality
	}{
		// ===== SDTV =====
		{"S07E23 .avi", SDTV},
		{"The.Series.S01E13.x264-CtrlSD", SDTV},
		{"The Series S02E01 HDTV XviD 2HD", SDTV},
		{"The Series S05E11 PROPER HDTV XviD 2HD", SDTV},
		{"The Series Show S02E08 HDTV x264 FTP", SDTV},
		{"The.Series.2011.S02E01.WS.PDTV.x264-TLA", SDTV},
		{"The.Series.2011.S02E01.WS.PDTV.x264-REPACK-TLA", SDTV},
		{"The Series S01E04 DSR x264 2HD", SDTV},
		{"The Series S01E04 Series Death Train DSR x264 MiNDTHEGAP", SDTV},
		{"The Series S11E03 has no periods or extension HDTV", SDTV},
		{"The.Series.S04E05.HDTV.XviD-LOL", SDTV},
		{"The.Series.S02E15.avi", SDTV},
		{"The.Series.S02E15.xvid", SDTV},
		{"The.Series.S02E15.divx", SDTV},
		{"The.Series.S03E06.HDTV-WiDE", SDTV},
		{"Series.S10E27.WS.DSR.XviD-2HD", SDTV},
		{"[HorribleSubs] The Series - 32 [480p]", SDTV},
		{"[CR] The Series - 004 [480p][48CE2D0F]", SDTV},
		{"[Hatsuyuki] The Series - 363 [848x480][ADE35E38]", SDTV},
		{"The.Series.S03.TVRip.XviD-NOGRP", SDTV},
		{"[HorribleSubs] The Series - 03 [360p].mkv", SDTV},

		// 540p without source returns Unknown (Sonarr parity)
		{"[SubsPlease] Series Title (540p) [AB649D32].mkv", Unknown},
		{"[Erai-raws] Series Title [540p][Multiple Subtitle].mkv", Unknown},

		// ===== DVD =====
		{"The.Series.S01E13.NTSC.x264-CtrlSD", DVD},
		{"The.Series.S03E06.DVDRip.XviD-WiDE", DVD},
		{"The.Series.S03E06.DVD.Rip.XviD-WiDE", DVD},
		{"the.Series.1x13.circles.ws.xvidvd-tns", DVD},
		{"the_Series.9x18.sunshine_days.ac3.ws_dvdrip_xvid-fov.avi", DVD},
		{"[FroZen] Series - 23 [DVD][7F6170E6]", DVD},
		{"[AniDL] Series - 26 -[360p][DVD][D - A][Exiled - Destiny]", DVD},

		// ===== WEBDL-480p =====
		{"The.Series.S01E10.The.Leviathan.480p.WEB-DL.x264-mSD", WEBDL480p},
		{"The.Series.S04E10.Glee.Actually.480p.WEB-DL.x264-mSD", WEBDL480p},
		{"The.SeriesS06E11.The.Santa.Simulation.480p.WEB-DL.x264-mSD", WEBDL480p},
		{"The.Series.S02E04.480p.WEB.DL.nSD.x264-NhaNc3", WEBDL480p},
		{"The.Series.S01E08.Das.geloeschte.Ich.German.Dubbed.DL.AmazonHD.x264-TVS", WEBDL480p},
		{"The.Series.S01E04.Rod.Trip.mit.meinem.Onkel.German.DL.NetflixUHD.x264", WEBDL480p},
		{"[HorribleSubs] Series Title! S01 [Web][MKV][h264][480p][AAC 2.0][Softsubs (HorribleSubs)]", WEBDL480p},
		{"Series.Title.S13E11.Ausgebacken.German.AmazonSD.h264-4SF", WEBDL480p},

		// ===== Bluray-480p =====
		{"SERIES.S03E01-06.DUAL.XviD.Bluray.AC3-REPACK.-HELLYWOOD.avi", Bluray480p},
		{"SERIES.S03E01-06.DUAL.BDRip.XviD.AC3.-HELLYWOOD", Bluray480p},
		{"SERIES.S03E01-06.DUAL.BDRip.X-viD.AC3.-HELLYWOOD", Bluray480p},
		{"SERIES.S03E01-06.DUAL.BDRip.AC3.-HELLYWOOD", Bluray480p},
		{"SERIES.S03E01-06.DUAL.BDRip.XviD.AC3.-HELLYWOOD.avi", Bluray480p},
		{"SERIES.S03E01-06.DUAL.XviD.Bluray.AC3.-HELLYWOOD.avi", Bluray480p},
		{"The.Series.S01E05.480p.BluRay.DD5.1.x264-HiSD", Bluray480p},
		{"The Series (BD)(640x480(RAW) (BATCH 1) (1-13)", Bluray480p},
		{"[Doki] Series - 02 (848x480 XviD BD MP3) [95360783]", Bluray480p},
		{"Adventures.of.Sonic.the.Hedgehog.S01.BluRay.480i.DD.2.0.AVC.REMUX-FraMeSToR", Bluray480p},
		{"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.480i.DD.2.0.AVC.REMUX-FraMeSToR", Bluray480p},

		// ===== WEBRip-480p =====
		{"The.Series.S02E10.480p.HULU.WEBRip.x264-Puffin", WEBRip480p},
		{"The.Series.S10E14.Techs.And.Balances.480p.AE.WEBRip.AAC2.0.x264-SEA", WEBRip480p},
		{"Series.Title.1x04.ITA.WEBMux.x264-NovaRip", WEBRip480p},

		// ===== Bluray-576p =====
		{"The.Series.S01E05.576p.BluRay.DD5.1.x264-HiSD", Bluray576p},

		// ===== HDTV-720p =====
		{"Series - S01E01 - Title [HDTV]", HDTV720p},
		{"Series - S01E01 - Title [HDTV-720p]", HDTV720p},
		{"The Series S04E87 REPACK 720p HDTV x264 aAF", HDTV720p},
		{"The.Series.S02E15.720p", HDTV720p},
		{"S07E23 - [HDTV-720p].mkv", HDTV720p},
		{"Series - S22E03 - MoneyBART - HD TV.mkv", HDTV720p},
		{"S07E23.mkv", HDTV720p},
		{"The.Series.S08E05.720p.HDTV.X264-DIMENSION", HDTV720p},
		{"The.Series.S02E15.mkv", HDTV720p},
		{"The.Series.S01E08.Tourmaline.Nepal.720p.HDTV.x264-DHD", HDTV720p},
		{"[Underwater-FFF] The Series - 01 (720p) [27AAA0A0]", HDTV720p},
		{"[Doki] The Series - 07 (1280x720 Hi10P AAC) [80AF7DDE]", HDTV720p},
		{"[Doremi].The.Series.5.Go.Go!.31.[1280x720].[C65D4B1F].mkv", HDTV720p},
		{"[HorribleSubs]_Series_Title_-_145_[720p]", HDTV720p},
		{"[Eveyuu] Series Title - 10 [Hi10P 1280x720 H264][10B23BD8]", HDTV720p},
		{"The.Series.US.S12E17.HR.WS.PDTV.X264-DIMENSION", HDTV720p},
		{"The Series S01E07 - Motor zmen (CZ)[TvRip][HEVC][720p]", HDTV720p},
		{"The.Series.S05E06.720p.HDTV.x264-FHD", HDTV720p},
		{"Series.Title.1x01.ITA.720p.x264-RlsGrp [01/54] - \"series.title.1x01.ita.720p.x264-rlsgrp.nfo\"", HDTV720p},
		{"[TMS-Remux].Series.Title.X.21.720p.[76EA1C53].mkv", HDTV720p},

		// ===== HDTV-1080p =====
		{"Under the Series S01E10 Let the Sonarr Begin 1080p", HDTV1080p},
		{"Series.S07E01.ARE.YOU.1080P.HDTV.X264-QCF", HDTV1080p},
		{"Series.S07E01.ARE.YOU.1080P.HDTV.x264-QCF", HDTV1080p},
		{"Series.S07E01.ARE.YOU.1080P.HDTV.proper.X264-QCF", HDTV1080p},
		{"Series - S01E01 - Title [HDTV-1080p]", HDTV1080p},
		{"[HorribleSubs] Series Title - 32 [1080p]", HDTV1080p},
		{"Series S01E07 - Sonarr zmen (CZ)[TvRip][HEVC][1080p]", HDTV1080p},
		{"The Online Series Alicization 04 vostfr FHD", HDTV1080p},
		{"Series Slayer 04 vostfr FHD.mkv", HDTV1080p},
		{"[Onii-ChanSub] The.Series - 02 vostfr (FHD 1080p 10bits).mkv", HDTV1080p},
		{"[Miaou] Series Title 02 VOSTFR FHD 10 bits", HDTV1080p},
		{"[mhastream.com]_Episode_05_FHD.mp4", HDTV1080p},
		{"[Kousei]_One_Series_ - _609_[FHD][648A87C7].mp4", HDTV1080p},
		{"Series culpable 1x02 Culpabilidad [HDTV 1080i AVC MP2 2.0 Sub][GrupoHDS]", HDTV1080p},
		{"Series como paso - 19x15 [344] Cuarenta anos de baile [HDTV 1080i AVC MP2 2.0 Sub][GrupoHDS]", HDTV1080p},
		{"Super.Seires.Go.S01E02.Depths.of.Sonarr.1080i.HDTV.DD5.1.H.264-NOGRP", HDTV1080p},

		// ===== HDTV-2160p =====
		{"My Title - S01E01 - EpTitle [HEVC 4k DTSHD-MA-6ch]", HDTV2160p},
		{"My Title - S01E01 - EpTitle [HEVC-4k DTSHD-MA-6ch]", HDTV2160p},
		{"My Title - S01E01 - EpTitle [4k HEVC DTSHD-MA-6ch]", HDTV2160p},

		// ===== WEBDL-720p =====
		{"Series S01E04 Mexicos Death Train 720p WEB DL", WEBDL720p},
		{"Series Five 0 S02E21 720p WEB DL DD5 1 H 264", WEBDL720p},
		{"Series S04E22 720p WEB DL DD5 1 H 264 NFHD", WEBDL720p},
		{"Series - S11E06 - D-Yikes! - 720p WEB-DL.mkv", WEBDL720p},
		{"The.Series.S02E15.720p.WEB-DL.DD5.1.H.264-SURFER", WEBDL720p},
		{"S07E23 - [WEBDL].mkv", WEBDL720p},
		{"Series S04E22 720p WEB-DL DD5.1 H264-EbP.mkv", WEBDL720p},
		{"Series.S04.720p.Web-Dl.Dd5.1.h264-P2PACK", WEBDL720p},
		{"Da.Series.Shows.S02E04.720p.WEB.DL.nSD.x264-NhaNc3", WEBDL720p},
		{"Series.Miami.S04E25.720p.iTunesHD.AVC-TVS", WEBDL720p},
		{"Series.S06E23.720p.WebHD.h264-euHD", WEBDL720p},
		{"Series.Title.2016.03.14.720p.WEB.x264-spamTV", WEBDL720p},
		{"Series.Title.2016.03.14.720p.WEB.h264-spamTV", WEBDL720p},
		{"Series.S01E08.Das.geloeschte.Ich.German.DD51.Dubbed.DL.720p.AmazonHD.x264-TVS", WEBDL720p},
		{"Series.Polo.S01E11.One.Hundred.Sonarrs.2015.German.DD51.DL.720p.NetflixUHD.x264.NewUp.by.Wunschtante", WEBDL720p},
		{"Series 2016 German DD51 DL 720p NetflixHD x264-TVS", WEBDL720p},
		{"Series.6x10.Basic.Sonarr.Repair.and.Replace.ITA.ENG.720p.WEB-DLMux.H.264-GiuseppeTnT", WEBDL720p},
		{"Series.6x11.Modern.Spy.ITA.ENG.720p.WEB.DLMux.H.264-GiuseppeTnT", WEBDL720p},
		{"The Series Was Dead 2010 S09E13 [MKV / H.264 / AC3/AAC / WEB / Dual Audio / Ingles / 720p]", WEBDL720p},
		{"into.the.Series.s03e16.h264.720p-web-handbrake.mkv", WEBDL720p},
		{"Series.S01E01.The.Sonarr.Principle.720p.WEB-DL.DD5.1.H.264-BD", WEBDL720p},
		{"Series.S03E05.Griebnitzsee.German.720p.MaxdomeHD.AVC-TVS", WEBDL720p},
		{"[HorribleSubs] Series Title! S01 [Web][MKV][h264][720p][AAC 2.0][Softsubs (HorribleSubs)]", WEBDL720p},
		{"[HorribleSubs] Series Title! S01 [Web][MKV][h264][AAC 2.0][Softsubs (HorribleSubs)]", WEBDL720p},
		{"Series.Title.S04E13.960p.WEB-DL.AAC2.0.H.264-squalor", WEBDL720p},
		{"Series.Title.S16.DP.WEB.720p.DDP.5.1.H.264.PLEX", WEBDL720p},
		{"Series.Title.S01E01.Erste.Begegnungen.German.DD51.Synced.DL.720p.HBOMaxHD.AVC-TVS", WEBDL720p},
		{"Series.Title.S01E05.Tavora.greift.an.German.DL.720p.DisneyHD.h264-4SF", WEBDL720p},

		// ===== WEBRip-720p =====
		{"Series.Title.S04E01.720p.WEBRip.AAC2.0.x264-NFRiP", WEBRip720p},
		{"Series.Title.S01E07.A.Prayer.For.Mad.Sweeney.720p.AMZN.WEBRip.DD5.1.x264-NTb", WEBRip720p},
		{"Series.Title.S07E01.A.New.Home.720p.DSNY.WEBRip.AAC2.0.x264-TVSmash", WEBRip720p},
		{"Series.Title.1x04.ITA.720p.WEBMux.x264-NovaRip", WEBRip720p},

		// ===== WEBDL-1080p =====
		{"Series S09E03 1080p WEB DL DD5 1 H264 NFHD", WEBDL1080p},
		{"Two and a Half Developers of the Series S10E03 1080p WEB DL DD5 1 H 264 NFHD", WEBDL1080p},
		{"Series.S08E01.1080p.WEB-DL.DD5.1.H264-NFHD", WEBDL1080p},
		{"Its.Always.Sonarrs.Fault.S08E01.1080p.WEB-DL.proper.AAC2.0.H.264", WEBDL1080p},
		{"This is an Easter Egg S10E03 1080p WEB DL DD5 1 H 264 REPACK NFHD", WEBDL1080p},
		{"Series.S04E09.Swan.Song.1080p.WEB-DL.DD5.1.H.264-ECI", WEBDL1080p},
		{"The.Big.Easter.Theory.S06E11.The.Sonarr.Simulation.1080p.WEB-DL.DD5.1.H.264", WEBDL1080p},
		{"Sonarr's.Baby.S01E02.Night.2.[WEBDL-1080p].mkv", WEBDL1080p},
		{"Series.Title.2016.03.14.1080p.WEB.x264-spamTV", WEBDL1080p},
		{"Series.Title.2016.03.14.1080p.WEB.h264-spamTV", WEBDL1080p},
		{"Series.S01.1080p.WEB-DL.AAC2.0.AVC-TrollHD", WEBDL1080p},
		{"Series Title S06E08 1080p WEB h264-EXCLUSIVE", WEBDL1080p},
		{"Series Title S06E08 No One PROPER 1080p WEB DD5 1 H 264-EXCLUSIVE", WEBDL1080p},
		{"Series Title S06E08 No One PROPER 1080p WEB H 264-EXCLUSIVE", WEBDL1080p},
		{"The.Series.S25E21.Pay.No1.1080p.WEB-DL.DD5.1.H.264-NTb", WEBDL1080p},
		{"Series.S01E08.Das.geloeschte.Ich.German.DD51.Dubbed.DL.1080p.AmazonHD.x264-TVS", WEBDL1080p},
		{"Death.Series.2017.German.DD51.DL.1080p.NetflixHD.x264-TVS", WEBDL1080p},
		{"Series.S01E08.Pro.Gamer.1440p.BKPL.WEB-DL.H.264-LiGHT", WEBDL1080p},
		{"Series.Title.S04E11.Teddy's.Choice.FHD.1080p.Web-DL", WEBDL1080p},
		{"Series.S04E03.The.False.Bride.1080p.NF.WEB.DDP5.1.x264-NTb[rartv]", WEBDL1080p},
		{"Series.Title.S02E02.This.Year.Will.Be.Different.1080p.AMZN.WEB...", WEBDL1080p},
		{"Series.Title.S02E02.This.Year.Will.Be.Different.1080p.AMZN.WEB.", WEBDL1080p},
		{"Series Title - S01E11 2020 1080p Viva MKV WEB", WEBDL1080p},
		{"[HorribleSubs] Series Title! S01 [Web][MKV][h264][1080p][AAC 2.0][Softsubs (HorribleSubs)]", WEBDL1080p},
		{"[LostYears] Series Title - 01-17 (WEB 1080p x264 10-bit AAC) [Dual-Audio]", WEBDL1080p},
		{"Series.and.Titles.S01.1080p.NF.WEB.DD2.0.x264-SNEAkY", WEBDL1080p},
		{"Series.Title.S02E02.This.Year.Will.Be.Different.1080p.WEB.H 265", WEBDL1080p},
		{"Series Title Season 2 [WEB 1080p HEVC Opus] [Netaro]", WEBDL1080p},
		{"Series Title Season 2 (WEB 1080p HEVC Opus) [Netaro]", WEBDL1080p},
		{"Series.Title.S01E01.Erste.Begegnungen.German.DD51.Synced.DL.1080p.HBOMaxHD.AVC-TVS", WEBDL1080p},
		{"Series.Title.S01E05.Tavora.greift.an.German.DL.1080p.DisneyHD.h264-4SF", WEBDL1080p},
		{"Series.Title.S02E04.German.Dubbed.DL.AAC.1080p.WEB.AVC-GROUP", WEBDL1080p},

		// ===== WEBRip-1080p =====
		{"Series.Title.S04E01.iNTERNAL.1080p.WEBRip.x264-QRUS", WEBRip1080p},
		{"Series.Title.S07E20.1080p.AMZN.WEBRip.DDP5.1.x264-ViSUM ac3.(NLsub)", WEBRip1080p},
		{"Series.Title.S03E09.1080p.NF.WEBRip.DD5.1.x264-ViSUM", WEBRip1080p},
		{"The Series 42 S09E13 1.54 GB WEB-RIP 1080p Dual-Audio 2019 MKV", WEBRip1080p},
		{"Series.Title.1x04.ITA.1080p.WEBMux.x264-NovaRip", WEBRip1080p},
		{"Series.Title.2019.S02E07.Chapter.15.The.Believer.4Kto1080p.DSNYP.Webrip.x265.10bit.EAC3.5.1.Atmos.GokiTAoE", WEBRip1080p},
		{"Series.Title.S01.1080p.AMZN.WEB-Rip.DDP5.1.H.264-Telly", WEBRip1080p},

		// ===== WEBDL-2160p =====
		{"Series.Title.2016.03.14.2160p.WEB.x264-spamTV", WEBDL2160p},
		{"Series.Title.2016.03.14.2160p.WEB.h264-spamTV", WEBDL2160p},
		{"Series.Title.2016.03.14.2160p.WEB.PROPER.h264-spamTV", WEBDL2160p},
		{"House.of.Sonarr.AK.s05e13.4K.UHD.WEB.DL", WEBDL2160p},
		{"House.of.Sonarr.AK.s05e13.UHD.4K.WEB.DL", WEBDL2160p},
		{"[HorribleSubs] Series Title! S01 [Web][MKV][h264][2160p][AAC 2.0][Softsubs (HorribleSubs)]", WEBDL2160p},
		{"Series Title S02 2013 WEB-DL 4k H265 AAC 2Audio-HDSWEB", WEBDL2160p},
		{"Series.Title.S02E02.This.Year.Will.Be.Different.2160p.WEB.H.265", WEBDL2160p},
		{"Series.Title.S02E04.German.Dubbed.DL.AAC.2160p.DV.HDR.WEB.HEVC-GROUP", WEBDL2160p},

		// ===== WEBRip-2160p =====
		{"Series S01E01.2160P AMZN WEBRIP DD2.0 HI10P X264-TROLLUHD", WEBRip2160p},
		{"JUST ADD SONARR S01E01.2160P AMZN WEBRIP DD2.0 X264-TROLLUHD", WEBRip2160p},
		{"The.Man.In.The.Series.S01E01.2160p.AMZN.WEBRip.DD2.0.Hi10p.X264-TrollUHD", WEBRip2160p},
		{"The Man In the Series S01E01 2160p AMZN WEBRip DD2.0 Hi10P x264-TrollUHD", WEBRip2160p},
		{"House.of.Sonarr.AK.S05E08.Chapter.60.2160p.NF.WEBRip.DD5.1.x264-NTb.NLsubs", WEBRip2160p},
		{"Sonarr Saves the World S01 2160p Netflix WEBRip DD5.1 x264-TrollUHD", WEBRip2160p},

		// ===== Bluray-720p =====
		{"SERIES.S03E01-06.DUAL.Bluray.AC3.-HELLYWOOD.avi", Bluray720p},
		{"Series - S01E03 - Come Fly With Me - 720p BluRay.mkv", Bluray720p},
		{"The Big Series.S03E01.The Sonarr Can Opener.m2ts", Bluray720p},
		{"Series.S01E02.Chained.Sonarr.[Bluray720p].mkv", Bluray720p},
		{"[FFF] DATE A Sonarr Dev - 01 [BD][720p-AAC][0601BED4]", Bluray720p},
		{"[RandomRemux] Series - 01 [720p BD][043EA407].mkv", Bluray720p},
		{"[Kaylith] Series Friends Specials - 01 [BD 720p AAC][B7EEE164].mkv", Bluray720p},
		{"SERIES.S03E01-06.DUAL.Blu-ray.AC3.-HELLYWOOD.avi", Bluray720p},
		{"SERIES.S03E01-06.DUAL.720p.Blu-ray.AC3.-HELLYWOOD.avi", Bluray720p},
		{"[Elysium]Lucky.Series.01(BD.720p.AAC.DA)[0BB96AD8].mkv", Bluray720p},
		{"Series.Galaxy.S01E01.33.720p.HDDVD.x264-SiNNERS.mkv", Bluray720p},
		{"The.Series.S01E07.RERIP.720p.BluRay.x264-DEMAND", Bluray720p},
		{"Series.Black.1x01.Selezione.Naturale.ITA.720p.BDMux.x264-NovaRip", Bluray720p},
		{"Series.Hunter.S02.720p.Blu-ray.Remux.AVC.FLAC.2.0-SiCFoI", Bluray720p},
		{"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.720p.DD.2.0.AVC.REMUX-FraMeSToR", Bluray720p},

		// ===== Bluray-1080p =====
		{"Series - S01E03 - Come Fly With Me - 1080p BluRay.mkv", Bluray1080p},
		{"Sonarr.Of.Series.S02E13.1080p.BluRay.x264-AVCDVD", Bluray1080p},
		{"Series.S01E02.Chained.Heat.[Bluray1080p].mkv", Bluray1080p},
		{"[FFF] Series no Muromi-san - 10 [BD][1080p-FLAC][0C4091AF]", Bluray1080p},
		{"[Kaylith] Series Friends Specials - 01 [BD 1080p FLAC][429FD8C7].mkv", Bluray1080p},
		{"[Zurako] Log Series - 01 - The Sonarr (BD 1080p AAC) [7AE12174].mkv", Bluray1080p},
		{"SERIES.S03E01-06.DUAL.1080p.Blu-ray.AC3.-HELLYWOOD.avi", Bluray1080p},
		{"[Coalgirls]_Series!!_01_(1920x1080_Blu-ray_FLAC)_[8370CB8F].mkv", Bluray1080p},
		{"Planet.Series.S01E11.Code.Deep.1080p.HD-DVD.DD.VC1-TRB", Bluray1080p},
		{"S for Series 2005 1080p UHD BluRay DD+7.1 x264-LoRD.mkv", Bluray1080p},
		{"Series.Title.2011.1080p.UHD.BluRay.DD5.1.HDR.x265-CtrlHD.mkv", Bluray1080p},
		{"Fall.Of.The.Release.Groups.S02E13.1080p.BDLight.x265-AVCDVD", Bluray1080p},

		// ===== Bluray-1080p Remux =====
		{"Series!!! on ICE - S01E12[JP BD Remux][ENG subs]", Bluray1080pRemux},
		{"Series.Title.S01E08.The.Well.BluRay.1080p.AVC.DTS-HD.MA.5.1.REMUX-FraMeSToR", Bluray1080pRemux},
		{"Series.Title.2x11.Nato.Per.La.Truffa.Bluray.Remux.AVC.1080p.AC3.ITA", Bluray1080pRemux},
		{"Series.Title.2x11.Nato.Per.La.Truffa.Bluray.Remux.AVC.AC3.ITA", Bluray1080pRemux},
		{"Series.Title.S03E01.The.Calm.1080p.DTS-HD.MA.5.1.AVC.REMUX-FraMeSToR", Bluray1080pRemux},
		{"Series Title Season 2 (BDRemux 1080p HEVC FLAC) [Netaro]", Bluray1080pRemux},
		{"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.1080p.DD.2.0.AVC.REMUX-FraMeSToR", Bluray1080pRemux},
		{"Series Title S01 2018 1080p BluRay Hybrid-REMUX AVC TRUEHD 5.1 Dual Audio-ZR-", Bluray1080pRemux},
		{"Series.Title.S01.2018.1080p.BluRay.Hybrid-REMUX.AVC.TRUEHD.5.1.Dual.Audio-ZR-", Bluray1080pRemux},

		// ===== Bluray-2160p =====
		{"Series.Title.US.s05e13.4K.UHD.Bluray", Bluray2160p},
		{"Series.Title.US.s05e13.UHD.4K.Bluray", Bluray2160p},
		{"[DameDesuYo] Series Bundle - Part 1 (BD 4K 8bit FLAC)", Bluray2160p},
		{"Series.Title.2014.2160p.UHD.BluRay.X265-IAMABLE.mkv", Bluray2160p},

		// ===== Bluray-2160p Remux =====
		{"Series!!! on ICE - S01E12[JP BD 2160p Remux][ENG subs]", Bluray2160pRemux},
		{"Series.Title.S01E08.The.Sonarr.BluRay.2160p.AVC.DTS-HD.MA.5.1.REMUX-FraMeSToR", Bluray2160pRemux},
		{"Series.Title.2x11.Nato.Per.The.Sonarr.Bluray.Remux.AVC.2160p.AC3.ITA", Bluray2160pRemux},
		{"[Dolby Vision] Sonarr.of.Series.S07.MULTi.UHD.BLURAY.REMUX.DV-NoTag", Bluray2160pRemux},
		{"Adventures.of.Sonic.the.Hedgehog.S01E01.Best.Hedgehog.2160p.DD.2.0.AVC.REMUX-FraMeSToR", Bluray2160pRemux},
		{"Series Title S01 2018 2160p BluRay Hybrid-REMUX AVC TRUEHD 5.1 Dual Audio-ZR-", Bluray2160pRemux},
		{"Series.Title.S01.2018.2160p.BluRay.Hybrid-REMUX.AVC.TRUEHD.5.1.Dual.Audio-ZR-", Bluray2160pRemux},

		// ===== Raw-HD =====
		{"POI S02E11 1080i HDTV DD5.1 MPEG2-TrollHD", RAWHD},
		{"How I Met Your Developer S01E18 Nothing Good Happens After Sonarr 720p HDTV DD5.1 MPEG2-TrollHD", RAWHD},
		{"The Series S01E11 The Finals 1080i HDTV DD5.1 MPEG2-TrollHD", RAWHD},
		{"Series.Title.S07E11.1080i.HDTV.DD5.1.MPEG2-NTb.ts", RAWHD},
		{"Game of Series S04E10 1080i HDTV MPEG2 DD5.1-CtrlHD.ts", RAWHD},
		{"Series.Title.S02E05.1080i.HDTV.DD2.0.MPEG2-NTb.ts", RAWHD},
		{"Show - S03E01 - Episode Title Raw-HD.ts", RAWHD},
		{"Series.Title.S10E09.Title.1080i.UPSCALE.HDTV.DD5.1.MPEG2-zebra", RAWHD},
		{"Series.Title.2011-08-04.1080i.HDTV.MPEG-2-CtrlHD", RAWHD},

		// ===== Additional test cases =====
		// WEB-DL variants
		{"The.Show.S01E01.1080p.WEB-DL.DD5.1.H.264-GROUP", WEBDL1080p},
		{"The.Show.S01E02.720p.WEB-DL.AAC2.0.H.264-GROUP", WEBDL720p},
		{"The.Show.S01E03.2160p.WEB-DL.DDP5.1.H.265-GROUP", WEBDL2160p},
		{"The.Show.S01E04.480p.WEB-DL.AAC2.0.H.264-GROUP", WEBDL480p},

		// WEBRip variants
		{"The.Show.S02E01.1080p.WEBRip.x264-GROUP", WEBRip1080p},
		{"The.Show.S02E02.720p.WEBRip.AAC2.0.H.264-GROUP", WEBRip720p},
		{"The.Show.S02E03.2160p.WEBRip.x265-GROUP", WEBRip2160p},
		{"The.Show.S02E04.480p.WEBRip.x264-GROUP", WEBRip480p},

		// HDTV variants
		{"The.Show.S03E01.1080p.HDTV.x264-GROUP", HDTV1080p},
		{"The.Show.S03E02.720p.HDTV.x264-GROUP", HDTV720p},
		{"The.Show.S03E03.2160p.HDTV.x265-GROUP", HDTV2160p},
		{"The.Show.S03E04.HDTV.x264-GROUP", SDTV},

		// BluRay variants
		{"The.Show.S04E01.1080p.BluRay.x264-GROUP", Bluray1080p},
		{"The.Show.S04E02.720p.BluRay.x264-GROUP", Bluray720p},
		{"The.Show.S04E03.2160p.BluRay.x265-GROUP", Bluray2160p},
		{"The.Show.S04E04.480p.BluRay.x264-GROUP", Bluray480p},
		{"The.Show.S04E05.576p.BluRay.x264-GROUP", Bluray576p},

		// REMUX variants
		{"The.Show.S05E01.1080p.BluRay.REMUX.AVC.DTS-HD.MA.5.1-GROUP", Bluray1080pRemux},
		{"The.Show.S05E02.2160p.UHD.BluRay.REMUX.HDR.HEVC.DTS-HD.MA.5.1-GROUP", Bluray2160pRemux},

		// DVD
		{"The.Show.S06E01.DVDRip.x264-GROUP", DVD},
		{"The.Show.S06E02.DVD.x264-GROUP", DVD},

		// BDRip/BRRip variants
		{"The.Show.S01E01.1080p.BDRip.x264-GROUP", Bluray1080p},
		{"The.Show.S01E02.720p.BDRip.x264-GROUP", Bluray720p},
		{"The.Show.S01E03.1080p.BRRip.x264-GROUP", Bluray1080p},
		{"The.Show.S01E04.720p.BRRip.x264-GROUP", Bluray720p},
		{"The.Show.S01E05.BDRip.x264-GROUP", Bluray480p},

		// PDTV/DSR/TVRip/SDTV variants
		{"The.Show.S01E01.PDTV.x264-GROUP", SDTV},
		{"The.Show.S01E02.720p.PDTV.x264-GROUP", HDTV720p},
		{"The.Show.S01E03.1080p.PDTV.x264-GROUP", HDTV1080p},
		{"The.Show.S01E04.DSR.x264-GROUP", SDTV},
		{"The.Show.S01E05.WS.DSR.x264-GROUP", SDTV},
		{"The.Show.S01E06.TVRip.x264-GROUP", SDTV},
		{"The.Show.S01E07.SDTV.x264-GROUP", SDTV},

		// Alternative resolution detection
		{"The.Show.S01E01.UHD.BluRay.x265-GROUP", Bluray2160p},
		{"The.Show.S01E02.[4K].BluRay.x265-GROUP", Bluray2160p},
		{"The.Show.S01E03.UHD.WEB-DL.x265-GROUP", WEBDL2160p},

		// Anime patterns
		{"[SubGroup] Anime Title - 01 [BD 1080p]", Bluray1080p},
		{"[SubGroup] Anime Title - 02 [BD 720p]", Bluray720p},
		{"[SubGroup] Anime Title - 03 [BD 2160p]", Bluray2160p},
		{"[SubGroup] Anime Title - 04 [WEB 1080p]", WEBDL1080p},
		{"[SubGroup] Anime Title - 05 [WEB 720p]", WEBDL720p},
		{"[SubGroup] Anime Title - 06 (WEB 1080p)", WEBDL1080p},

		// BDLight variant
		{"The.Show.S01E01.1080p.BDLight.x264-GROUP", Bluray1080p},

		// Raw-HD variants
		{"The.Show.S01E01.1080i.RawHD.DD5.1-GROUP", RAWHD},
		{"The.Show.S01E02.Raw-HD.x264-GROUP", RAWHD},
		{"The.Show.S01E03.720p.HDTV.MPEG2-GROUP", RAWHD},

		// Streaming service specific patterns
		{"The.Show.S09E01.1080p.AMZN.WEB-DL.DDP5.1.H.264-GROUP", WEBDL1080p},
		{"The.Show.S09E02.1080p.NF.WEB-DL.DDP5.1.H.264-GROUP", WEBDL1080p},
		{"The.Show.S09E03.2160p.DSNP.WEB-DL.DDP5.1.H.265-GROUP", WEBDL2160p},

		// More edge cases
		{"The.Show.2024.S01E01.1080p.WEB.H264-GROUP", WEBDL1080p},
		{"The.Show.S01E01.1080i.HDTV.DD5.1.MPEG2-GROUP", RAWHD},
		{"The.Show.S01E01.4K.UHD.BluRay.x265-GROUP", Bluray2160p},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result := ParseQuality(tt.title)
			if result.Quality != tt.expected {
				t.Errorf("ParseQuality(%q) = %v (%s), want %v (%s)",
					tt.title, result.Quality, result.Quality.String(),
					tt.expected, tt.expected.String())
			}
		})
	}
}

func TestQualityModelMethods(t *testing.T) {
	qm := QualityModel{
		Quality: Bluray1080pRemux,
		Revision: Revision{
			Version:  2,
			IsRepack: true,
		},
	}

	if qm.Source() != "BluRay" {
		t.Errorf("QualityModel.Source() = %v, want %v", qm.Source(), "BluRay")
	}
	if qm.Resolution() != "1080p" {
		t.Errorf("QualityModel.Resolution() = %v, want %v", qm.Resolution(), "1080p")
	}
	if !qm.IsRemux() {
		t.Errorf("QualityModel.IsRemux() = %v, want %v", qm.IsRemux(), true)
	}
	if qm.String() != "Bluray-1080p Remux v2 [REPACK]" {
		t.Errorf("QualityModel.String() = %v, want %v", qm.String(), "Bluray-1080p Remux v2 [REPACK]")
	}
	if qm.Full() != "Bluray-1080p Remux" {
		t.Errorf("QualityModel.Full() = %v, want %v", qm.Full(), "Bluray-1080p Remux")
	}
	if qm.Version() != 2 {
		t.Errorf("QualityModel.Version() = %v, want %v", qm.Version(), 2)
	}
}

func TestFull(t *testing.T) {
	tests := []struct {
		name     string
		quality  Quality
		expected string
	}{
		{"HDTV720p", HDTV720p, "HDTV-720p"},
		{"WEBDL1080p", WEBDL1080p, "WEBDL-1080p"},
		{"Bluray2160pRemux", Bluray2160pRemux, "Bluray-2160p Remux"},
		{"SDTV", SDTV, "SDTV"},
		{"Unknown", Unknown, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qm := QualityModel{Quality: tt.quality}
			if got := qm.Full(); got != tt.expected {
				t.Errorf("QualityModel.Full() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestListFields(t *testing.T) {
	fields := ListFields()
	if len(fields) == 0 {
		t.Error("ListFields() returned empty slice")
	}

	// Check that expected fields are present
	fieldNames := make(map[string]bool)
	for _, field := range fields {
		fieldNames[field.Name] = true
	}

	expectedFields := []string{"Full", "Resolution", "Source", "IsRemux", "IsRepack", "Version"}
	for _, expected := range expectedFields {
		if !fieldNames[expected] {
			t.Errorf("ListFields() missing expected field: %s", expected)
		}
	}
}

func TestGetField(t *testing.T) {
	qm := QualityModel{
		Quality: HDTV720p,
		Revision: Revision{
			Version:  2,
			IsRepack: true,
		},
	}

	tests := []struct {
		name     string
		expected interface{}
	}{
		{"Full", "HDTV-720p"},
		{"Resolution", "720p"},
		{"Source", "HDTV"},
		{"IsRemux", false},
		{"IsRepack", true},
		{"Version", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetField(tt.name, qm)
			if err != nil {
				t.Errorf("GetField(%s) error = %v", tt.name, err)
				return
			}
			if got != tt.expected {
				t.Errorf("GetField(%s) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}

	// Test unknown field
	_, err := GetField("UnknownField", qm)
	if err == nil {
		t.Error("GetField() with unknown field should return error")
	}
}
