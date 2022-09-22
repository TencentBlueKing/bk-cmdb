/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <stdio.h>
#include <time.h>
#include <string.h>

#include "tools/strings.h"


#include "time.h"
// Date

namespace gse {
namespace data {
Date::Date()
{
	Init();
}

Date::Date(int year, int month, int day)
{

	m_year = year;
	m_month = month;
	m_day = day;
}

void Date::Init()
{

	time_t tt = time((time_t*) (NULL));
	struct tm* ptrTmp = gmtime(&tt);
	m_year = ptrTmp->tm_year + 1900;
	m_month = ptrTmp->tm_mon + 1;
	m_day = ptrTmp->tm_mday;
}

bool Date::operator ==(const Date& date)
{

	return m_year == date.m_year && m_month == date.m_month && m_day == date.m_day;
}

bool Date::operator !=(const Date& date)
{

	return !this->operator ==(date);
}

void Date::Parse(const char* date)
{

	m_year = m_month = m_day = 0;

	unsigned int len = strlen(date);
	for (unsigned int i = 0; i < len; ++i)
	{

		unsigned short tmp = 0;
		while (isdigit(date[i]))
		{

			tmp = tmp * 10 + date[i] - '0';
			++i;
		}

		if (m_year == 0)
		{

			m_year = tmp;
		}
		else if (m_month == 0)
		{

			m_month = tmp;
		}
		else if (m_day == 0)
		{

			m_day = tmp;
		}

	}
}

void Date::Year(int year)
{

	m_year = year;
}

int Date::Year() const
{

	return m_year;
}

void Date::Month(int month)
{

	m_month = month;
}

int Date::Month() const
{

	return m_month;
}

void Date::Day(int day)
{

	m_day = day;
}

int Date::Day() const
{

	return m_day;
}

std::string Date::ToString()
{

    memset(m_date, 0, DATE_LENGTH);
    gse::tools::strings::Snprintf(m_date, DATE_LENGTH, "%04d-%02d-%02d", m_year, m_month, m_day);

	return std::string(m_date);
}

// Time
Time::Time()
{

	Init();
}

Time::Time(int hour, int minute, int second)
{

	m_hour = hour;
	m_minute = minute;
	m_second = second;
}
void Time::Init()
{

	time_t tt = time((time_t*) (NULL));
	struct tm* ptrTmp = gmtime(&tt);

	m_hour = ptrTmp->tm_hour;
	m_minute = ptrTmp->tm_min;
	m_second = ptrTmp->tm_sec;
}

bool Time::operator ==(const Time& time)
{

	return m_second == time.m_second && m_minute == time.m_minute && m_hour == time.m_hour;
}

bool Time::operator !=(const Time& time)
{

	return !this->operator ==(time);
}

void Time::Parse(const char* time)
{

	m_hour = m_minute = m_second = 0;

	unsigned int len = strlen(time);
	for (unsigned int i = 0; i < len; ++i)
	{

		unsigned short tmp = 0;
		while (isdigit(time[i]))
		{

			tmp = tmp * 10 + time[i] - '0';
			++i;
		}

		if (m_hour == 0)
		{

			m_hour = tmp;
		}
		else if (m_minute == 0)
		{

			m_minute = tmp;
		}
		else if (m_second == 0)
		{

			m_second = tmp;
			return;
		}

	}
}

void Time::Hour(int hour)
{

	m_hour = hour;
}

int Time::Hour() const
{

	return m_hour;
}

void Time::Minute(int minute)
{

	m_minute = minute;
}

int Time::Minute() const
{

	return m_minute;
}

void Time::Second(int second)
{

	m_second = second;
}

int Time::Second() const
{

	return m_second;
}

std::string Time::ToString()
{

    memset(m_time, 0, TIME_LENGTH);
    gse::tools::strings::Snprintf(m_time, TIME_LENGTH, "%02d:%02d:%02d", m_hour, m_minute, m_second);
	return std::string(m_time);
}

// DateTime
DateTime::DateTime()
{

	Init();
}

DateTime::DateTime(const Date& date)
{

	Init();

	m_systemTime.tm_year = date.Year() - 1900;
	m_systemTime.tm_mon = date.Month() - 1;
	m_systemTime.tm_mday = date.Day();

	m_timeT = mktime(&m_systemTime);
}

DateTime::DateTime(const Time& time)
{

	Init();

	m_systemTime.tm_hour = time.Hour();
	m_systemTime.tm_min = time.Minute();
	m_systemTime.tm_sec = time.Second();

	m_timeT = mktime(&m_systemTime);

}

DateTime::DateTime(const Date& date, const Time& time)
{
	Init();

	m_systemTime.tm_year = date.Year() - 1900;
	m_systemTime.tm_mon = date.Month() - 1;

	m_systemTime.tm_mday = date.Day();
	m_systemTime.tm_hour = time.Hour();

	m_systemTime.tm_min = time.Minute();
	m_systemTime.tm_sec = time.Second();

	m_timeT = mktime(&m_systemTime);

}

DateTime::DateTime(time_t dateTime)
{

	m_timeT = dateTime;
	memcpy(&m_systemTime, gmtime(&m_timeT), sizeof(struct tm));
}

DateTime::DateTime(int year, int month, int day, int hour, int minute, int second)
{

	Init();

	m_systemTime.tm_year = year - 1900;
	m_systemTime.tm_mon = month - 1;
	m_systemTime.tm_mday = day;
	m_systemTime.tm_hour = hour;
	m_systemTime.tm_min = minute;
	m_systemTime.tm_sec = second;

	m_timeT = mktime(&m_systemTime);
}

void DateTime::Init()
{

	time(&m_timeT);
	memcpy(&m_systemTime, gmtime(&m_timeT), sizeof(struct tm));
}

time_t DateTime::UTCSecondsFrom1970()
{
    return this->ToUTCSeconds(); 
}

bool DateTime::operator==(const DateTime& time)
{

	return m_systemTime.tm_sec == time.m_systemTime.tm_sec && m_systemTime.tm_min == time.m_systemTime.tm_min && m_systemTime.tm_hour == time.m_systemTime.tm_hour
			&& m_systemTime.tm_mday == time.m_systemTime.tm_mday && m_systemTime.tm_mon == time.m_systemTime.tm_mon && m_systemTime.tm_year == time.m_systemTime.tm_year;
}

bool DateTime::operator!=(const DateTime& time)
{

	return !this->operator ==(time);
}

bool DateTime::operator>(const DateTime& time)
{

	return m_systemTime.tm_year > time.m_systemTime.tm_year || m_systemTime.tm_mon > time.m_systemTime.tm_mon || m_systemTime.tm_mday > time.m_systemTime.tm_mday
			|| m_systemTime.tm_hour > time.m_systemTime.tm_hour || m_systemTime.tm_min > time.m_systemTime.tm_min || m_systemTime.tm_sec > time.m_systemTime.tm_sec;
}

bool DateTime::operator>=(const DateTime& time)
{

	return m_systemTime.tm_year >= time.m_systemTime.tm_year || m_systemTime.tm_mon >= time.m_systemTime.tm_mon || m_systemTime.tm_mday >= time.m_systemTime.tm_mday
			|| m_systemTime.tm_hour >= time.m_systemTime.tm_hour || m_systemTime.tm_min >= time.m_systemTime.tm_min || m_systemTime.tm_sec >= time.m_systemTime.tm_sec;
}

bool DateTime::operator<(const DateTime& time)
{

	return m_systemTime.tm_year < time.m_systemTime.tm_year || m_systemTime.tm_mon < time.m_systemTime.tm_mon || m_systemTime.tm_mday < time.m_systemTime.tm_mday
			|| m_systemTime.tm_hour < time.m_systemTime.tm_hour || m_systemTime.tm_min < time.m_systemTime.tm_min || m_systemTime.tm_sec < time.m_systemTime.tm_sec;
}

bool DateTime::operator<=(const DateTime& time)
{

	return m_systemTime.tm_year <= time.m_systemTime.tm_year || m_systemTime.tm_mon <= time.m_systemTime.tm_mon || m_systemTime.tm_mday <= time.m_systemTime.tm_mday
			|| m_systemTime.tm_hour <= time.m_systemTime.tm_hour || m_systemTime.tm_min <= time.m_systemTime.tm_min || m_systemTime.tm_sec <= time.m_systemTime.tm_sec;
}

DateTime& DateTime::operator=(const Date& date)
{

	m_systemTime.tm_year = date.Year() - 1900;
	m_systemTime.tm_mon = date.Month() - 1;
	m_systemTime.tm_mday = date.Day();

	m_timeT = mktime(&m_systemTime);

	return *this;
}

DateTime& DateTime::operator=(const Time& time)
{

	m_systemTime.tm_hour = time.Hour();
	m_systemTime.tm_min = time.Minute();
	m_systemTime.tm_sec = time.Second();

	m_timeT = mktime(&m_systemTime);

	return *this;
}

void DateTime::Parse(const char* dateTime)
{

	Init();

	unsigned int len = strlen(dateTime);
	for (unsigned int i = 0; i < len; ++i)
	{

		unsigned short tmp = 0;
		while (isdigit(dateTime[i]))
		{

			tmp = tmp * 10 + dateTime[i] - '0';
			++i;
		}

		if (m_systemTime.tm_year == 0)
		{

			m_systemTime.tm_year = tmp - 1900;
		}
		else if (m_systemTime.tm_mon == 0)
		{

			m_systemTime.tm_mon = tmp - 1;
		}
		else if (m_systemTime.tm_mday == 0)
		{

			m_systemTime.tm_mday = tmp;
		}
		else if (m_systemTime.tm_hour == 0)
		{

			m_systemTime.tm_hour = tmp;
		}
		else if (m_systemTime.tm_min == 0)
		{

			m_systemTime.tm_min = tmp;
		}
		else if (m_systemTime.tm_sec == 0)
		{

			m_systemTime.tm_sec = tmp;

			m_timeT = mktime(&m_systemTime);
			return;

		}

	}
}

void DateTime::Year(int year)
{

	m_systemTime.tm_year = year - 1900;
	m_timeT = mktime(&m_systemTime);
}

int DateTime::Year() const
{

	return m_systemTime.tm_year + 1900;
}

void DateTime::Month(int month)
{

	m_systemTime.tm_mon = month - 1;
	m_timeT = mktime(&m_systemTime);
}

int DateTime::Month() const
{

	return m_systemTime.tm_mon + 1;
}

void DateTime::Day(int day)
{

	m_systemTime.tm_mday = day;
	m_timeT = mktime(&m_systemTime);
}

int DateTime::Day() const
{

	return m_systemTime.tm_mday;
}

void DateTime::Hour(int hour)
{

	m_systemTime.tm_hour = hour;
	m_timeT = mktime(&m_systemTime);
}

int DateTime::Hour() const
{

	return m_systemTime.tm_hour;
}

void DateTime::Minute(int minute)
{

	m_systemTime.tm_min = minute;
	m_timeT = mktime(&m_systemTime);
}

int DateTime::Minute() const
{

	return m_systemTime.tm_min;
}

void DateTime::Second(int second)
{

	m_systemTime.tm_sec = second;
	m_timeT = mktime(&m_systemTime);
}

int DateTime::Second() const
{

	return m_systemTime.tm_sec;
}

DateTime DateTime::ToLocalDateTime()
{

	struct tm* ptrTmp = localtime(&m_timeT);
	return DateTime(ptrTmp->tm_year + 1900, ptrTmp->tm_mon + 1, ptrTmp->tm_mday, ptrTmp->tm_hour, ptrTmp->tm_min, ptrTmp->tm_sec);
}

DateTime DateTime::ToUTCDateTime()
{

	struct tm* ptrTmp = gmtime(&m_timeT);
	return DateTime(ptrTmp->tm_year + 1900, ptrTmp->tm_mon + 1, ptrTmp->tm_mday, ptrTmp->tm_hour, ptrTmp->tm_min, ptrTmp->tm_sec);
}

Date DateTime::ToDate()
{

	return Date(m_systemTime.tm_year + 1900, m_systemTime.tm_mon + 1, m_systemTime.tm_mday);
}

Time DateTime::ToTime()
{

	return Time(m_systemTime.tm_hour, m_systemTime.tm_min, m_systemTime.tm_sec);
}

std::string DateTime::ToLocalTimeString()
{

    memset(m_bufTime, 0, BUF_LENGTH);
    struct tm* ptrTmp = localtime(&m_timeT);
    gse::tools::strings::Snprintf(m_bufTime, BUF_LENGTH, "%02d:%02d:%02d", ptrTmp->tm_hour, ptrTmp->tm_min, ptrTmp->tm_sec);
    return std::string(m_bufTime);
}

std::string DateTime::ToLocalDateString()
{

    memset(m_bufTime, 0, BUF_LENGTH);
    struct tm* ptrTmp = localtime(&m_timeT);
    gse::tools::strings::Snprintf(m_bufTime, BUF_LENGTH, "%04d-%02d-%02d", ptrTmp->tm_year + 1900, ptrTmp->tm_mon + 1, ptrTmp->tm_mday);
    return std::string(m_bufTime);
}

std::string DateTime::ToStandardTimestampString()
{
	memset(m_bufTime, 0, BUF_LENGTH);
	struct tm* ptrTmp = localtime(&m_timeT);
	// ISO 8601 formater
    gse::tools::strings::Snprintf(m_bufTime, BUF_LENGTH, "%04d-%02d-%02dT%02d:%02d:%02d%+03d:00", ptrTmp->tm_year + 1900, ptrTmp->tm_mon + 1, ptrTmp->tm_mday, ptrTmp->tm_hour, ptrTmp->tm_min, ptrTmp->tm_sec,GetTimezone());
	return std::string(m_bufTime);
}
std::string DateTime::ToLocalDateTimeString()
{
    memset(m_bufTime, 0, BUF_LENGTH);
    struct tm* ptrTmp = localtime(&m_timeT);
    gse::tools::strings::Snprintf(m_bufTime, BUF_LENGTH, "%04d-%02d-%02d %02d:%02d:%02d", ptrTmp->tm_year + 1900, ptrTmp->tm_mon + 1, ptrTmp->tm_mday, ptrTmp->tm_hour, ptrTmp->tm_min, ptrTmp->tm_sec);
    return std::string(m_bufTime);
}

std::string DateTime::ToUTCTimeString()
{

    memset(m_bufTime, 0, BUF_LENGTH);
    gse::tools::strings::Snprintf(m_bufTime, BUF_LENGTH, "%02d:%02d:%02d", m_systemTime.tm_hour, m_systemTime.tm_min, m_systemTime.tm_sec);
    return std::string(m_bufTime);
}

std::string DateTime::ToUTCDateString()
{

    memset(m_bufTime, 0, BUF_LENGTH);
    gse::tools::strings::Snprintf(m_bufTime, BUF_LENGTH, "%04d-%02d-%02d", m_systemTime.tm_year + 1900, m_systemTime.tm_mon + 1, m_systemTime.tm_mday);
    return std::string(m_bufTime);
}

std::string DateTime::ToUTCDateTimeString()
{
    memset(m_bufTime, 0, BUF_LENGTH);
    gse::tools::strings::Snprintf(m_bufTime, BUF_LENGTH, "%04d-%02d-%02d %02d:%02d:%02d", m_systemTime.tm_year + 1900, m_systemTime.tm_mon + 1, m_systemTime.tm_mday, m_systemTime.tm_hour, m_systemTime.tm_min, m_systemTime.tm_sec);
    return std::string(m_bufTime);
}

time_t DateTime::ToUTCSeconds()
{
	return m_timeT;
}
int  DateTime::GetTimezone(int * dst /*= NULL*/)
{
#ifdef __linux__
	time_t timeutc;
	// get utc time
	time(&timeutc);

	// get the local time
	struct tm tmlocal;
	localtime_r(&timeutc, &tmlocal);
	if (NULL != dst)
	{
		*dst = tmlocal.tm_isdst;
	}
    if( NULL != dst)
    {
        *dst = tmlocal.tm_isdst;
    }

	time_t timelocal;
	struct tm tmgmt;

	// change tm to time_t
	timelocal = mktime(&tmlocal);

	// change it to gmt tm
	gmtime_r(&timeutc, &tmgmt);

	int timezone = tmlocal.tm_hour - tmgmt.tm_hour;
	if (timezone < -12)
	{
		timezone += 24;
	}
	else if (timezone > 12)
	{
		timezone -= 24;
	}

	return timezone;
#else
    return 0;
#endif
}
}
}
