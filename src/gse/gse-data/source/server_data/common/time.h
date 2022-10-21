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

#ifndef __GSE_TIME_H__
#define __GSE_TIME_H__

#include <string>
#include <time.h>

namespace gse {
namespace data {
/**
 *\class  Date
 *\brief  日期
 */
class Date
{
public:
    Date();
    explicit Date(int year, int month, int day);

public:
    bool operator==(const Date& date);
    bool operator!=(const Date& date);

public:

    void Parse(const char* date);

    void Year(int year);
    int Year() const;

    void Month(int month);
    int Month() const;

    void Day(int day);
    int Day() const;

public:

    std::string ToString();

private:
    void Init();

private:
    enum
    {
        DATE_LENGTH = 128
    };

private:
    int m_year;
    int m_month;
    int m_day;

    char m_date[DATE_LENGTH];

};

/**
 *\class Time
 *\brief 时间
 */
class Time
{
public:
    Time();
    explicit Time(int hour, int minute, int second);

public:
    bool operator==(const Time& time);
    bool operator!=(const Time& time);

public:

    void Parse(const char* time);

    void Hour(int hour);
    int Hour() const;

    void Minute(int minute);
    int Minute() const;

    void Second(int second);
    int Second() const;

public:

    std::string ToString();

private:
    void Init();

private:
    enum
    {
        TIME_LENGTH = 128
    };

private:
    int m_hour;
    int m_minute;
    int m_second;

    char m_time[TIME_LENGTH];

};

/**
 *\class  Time
 *\brief  时间对象
 */
class DateTime
{
public:
    DateTime();
    explicit DateTime(const Date& date);
    explicit DateTime(const Time& time);
    explicit DateTime(const Date& date, const Time& time);
    explicit DateTime(time_t dateTime);
    explicit DateTime(int year, int month, int day, int hour, int minute, int second);

public:
    static time_t getUnixTimeStamp()
    {
        return time(NULL);
    }
    time_t UTCSecondsFrom1970();

public:
    // 运算符重载
    bool operator==(const DateTime& time);
    bool operator!=(const DateTime& time);
    bool operator>(const DateTime& time);
    bool operator>=(const DateTime& time);
    bool operator<(const DateTime& time);
    bool operator<=(const DateTime& time);

    DateTime& operator=(const Date& date);
    DateTime& operator=(const Time& time);

public:

    void Parse(const char* dateTime);

    void Year(int year);
    int Year() const;

    void Month(int month);
    int Month() const;

    void Day(int day);
    int Day() const;

    void Hour(int hour);
    int Hour() const;

    void Minute(int minute);
    int Minute() const;

    void Second(int second);
    int Second() const;

public:
    DateTime ToLocalDateTime();
    DateTime ToUTCDateTime();
    Date ToDate();
    Time ToTime();

    std::string ToLocalTimeString();
    std::string ToLocalDateString();
    std::string ToLocalDateTimeString();
    std::string ToStandardTimestampString();

    std::string ToUTCTimeString();
    std::string ToUTCDateString();
    std::string ToUTCDateTimeString();
    time_t ToUTCSeconds();
    int   GetTimezone(int * dst = NULL);

private:
    void Init();

private:
    enum
    {
        BUF_LENGTH = 128
    };
private:
    struct tm m_systemTime;
    time_t m_timeT;
    char m_bufTime[BUF_LENGTH];
};
}

}
#endif
