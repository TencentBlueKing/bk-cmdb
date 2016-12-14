<?php if (!defined('BASEPATH'))
    exit('No direct script access allowed');
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */


/**
 * 工具库
 */
class ExcelUtility {

    public static $errInfo = '';

    /**
     * 计算excel列名
     * @param num int 列数
     * @return arrray 读取的数据
     */
    public function getCol($num) {
        /*
         * 递归方式实现根据列数返回列的字母标识
        */
        $arr = array(0 => 'Z', 1 => 'A', 2 => 'B', 3 => 'C', 4 => 'D', 5 => 'E', 6 => 'F', 7 => 'G', 8 => 'H', 9 => 'I',
                     10 => 'J', 11 => 'K', 12 => 'L', 13 => 'M', 14 => 'N', 15 => 'O', 16 => 'P', 17 => 'Q', 18 => 'R',
                     19 => 'S', 20 => 'T', 21 => 'U', 22 => 'V', 23 => 'W', 24 => 'X', 25 => 'Y', 26 => 'Z');
        if ($num == 0) {
            return '';
        }
        return self::getCol((int)(($num - 1) / 26)) . $arr[$num % 26];
    }

    /**
     * @PHPExcel v1.8.0 读excel
     * @param file string文件名
     * @return arrray 读取的数据
     */
    public static function readExcel($file) {
        if (!file_exists($file)) {
            self::$errInfo = $file . '文件不存在';
            return false;
        }

        include 'PHPExcel-1.8.0/PHPExcel/IOFactory.php';
        $ext = pathinfo($file, PATHINFO_EXTENSION);
        if ($ext == 'xlsx' || $ext == 'xls') {
            $objPHPExcel = PHPExcel_IOFactory::load($file);
        } else if ($ext == 'csv') {
            $objReader = PHPExcel_IOFactory::createReader('CSV')
                ->setDelimiter(',')
                ->setInputEncoding('GBK') //不设置将导致中文列内容返回boolean(false)或乱码
                ->setEnclosure('"')
                ->setLineEnding("\r\n")
                ->setSheetIndex(0);
            $objPHPExcel = $objReader->load($file);
        }

        $sheetData = $objPHPExcel->getActiveSheet()->toArray(null, true, true, true);
        $header = array();
        if (count($sheetData) > 1) {
            foreach ($sheetData[1] as $_k => $_v) {
                $header[$_k] = $_v;
            }
        }
        $data = array();
        foreach ($sheetData as $_i => $_d) {
            if ($_i === 1) {
                continue;
            }
            foreach ($_d as $_k => $_v) {
                $data[$header[$_k]][] = $_v;
            }

        }
        return $data;
    }

    /**
     * 导出文件至excel
     * @param $headArr
     * @param $data
     * @param $fileName
     * @return bool
     */
    public static function exportToExcel($headArr, $data, $fileName) {
        include 'PHPExcel-1.8.0/PHPExcel/IOFactory.php';
        if (empty($data) || !is_array($data)) {
            self::$errInfo = '数据格式不正确';
            return false;
        }

        if (empty($fileName)) {
            self::$errInfo = $fileName . '文件不存在';
            return false;
        }

        $date = date("Y_m_d", time());
        $fileName .= "_" . $date . ".xlsx";

        /*创建新的PHPExcel对象*/
        $objPHPExcel = new PHPExcel();
        $objProps = $objPHPExcel->getProperties();

        /*设置表头*/
        $num = 1;
        $rowArr = array();
        foreach($headArr as $v) {
            $colum = self::getCol($num);
            $rowArr[$v] = $colum;
            $objPHPExcel->setActiveSheetIndex(0)->setCellValue($colum.'1', $v);
            $num++;
        }

        $rowIndex = 2;
        $objActSheet = $objPHPExcel->getActiveSheet();
        foreach($data as $key => $rows) {
            foreach($headArr as $keyName=>$value) {
                $j = $rowArr[$value];
                $objActSheet->setCellValue($j . $rowIndex, (isset($rows[$keyName]) ? $rows[$keyName] : ''));
            }
            $rowIndex++;
        }

        $fileName = iconv("utf-8", "gb2312", $fileName);
        /*重命名表*/
        $objPHPExcel->getActiveSheet()->setTitle('Simple');
        /*设置活动单指数到第一个表,所以Excel打开这是第一个表*/
        $objPHPExcel->setActiveSheetIndex(0);
        /*将输出重定向到一个客户端web浏览器*/
        header('Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet');
        header("Content-Disposition: attachment; filename=\"$fileName\"");
        header('Cache-Control: max-age=0');
        $objWriter = PHPExcel_IOFactory::createWriter($objPHPExcel, 'Excel2007');
        $objWriter->save('php://output');
    }
}
