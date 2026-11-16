import React, { useState, useEffect } from "react";
import { View, Text, TouchableOpacity, TextInput, Modal, ActivityIndicator, Alert, ScrollView } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import MapView from "react-native-maps";
import * as Location from "expo-location";
import { Ionicons } from "@expo/vector-icons";

import style from "./styled";
import { listarAnimais, registrarAvistamento } from "../../../../services/API";

const mensagemErro = (err) => {
  const msg = err?.message || "";
  if (msg.includes("Network request failed")) {
    return "Sem conexão com o servidor. Verifique a internet e tente de novo.";
  }
  return msg || "Não foi possível registrar o avistamento.";
};

export default () => {
  const navigation = useNavigation();

  const [animais, setAnimais] = useState([]);
  const [erroAnimais, setErroAnimais] = useState(false);
  const [animalId, setAnimalId] = useState(null);
  const [descricao, setDescricao] = useState("");
  
   
  const [coordenadas, setCoordenadas] = useState(null);
  const [mapaVisivel, setMapaVisivel] = useState(false);
  const [regiaoMapa, setRegiaoMapa] = useState(null);

 
  const [loading, setLoading] = useState(false);
  const [obtendoLocalizacao, setObtendoLocalizacao] = useState(false);

  const carregarAnimais = async () => {
    try {
      setErroAnimais(false);
      setAnimais(await listarAnimais());
    } catch (err) {
      console.error("Erro ao carregar animais:", err);
      setErroAnimais(true);
    }
  };

  useEffect(() => {
    carregarAnimais();
    solicitarPermissaoLocalizacao();
  }, []);

  const solicitarPermissaoLocalizacao = async () => {
    setObtendoLocalizacao(true);
    let { status } = await Location.requestForegroundPermissionsAsync();
    if (status !== 'granted') {
      Alert.alert("Permissão negada", "Precisamos de acesso à localização para marcar a ocorrência.");
      setObtendoLocalizacao(false);
      return;
    }

    let location = await Location.getCurrentPositionAsync({});
    const coordAtual = {
      latitude: location.coords.latitude,
      longitude: location.coords.longitude,
      latitudeDelta: 0.005,
      longitudeDelta: 0.005,
    };
    setRegiaoMapa(coordAtual);
    setObtendoLocalizacao(false);
  };

  const abrirMapa = () => {
    if (!regiaoMapa) {
      solicitarPermissaoLocalizacao();
    }
    setMapaVisivel(true);
  };

  const confirmarLocalizacao = () => {
    setCoordenadas({
      latitude: regiaoMapa.latitude,
      longitude: regiaoMapa.longitude,
    });
    setMapaVisivel(false);
  };

  const salvarOcorrencia = async () => {
    if (!animalId || !coordenadas) {
      Alert.alert("Atenção", "Escolha o animal e marque o local.");
      return;
    }

    try {
      setLoading(true);
      
      await registrarAvistamento({
        animal_id: animalId,
        latitude: coordenadas.latitude,
        longitude: coordenadas.longitude,
        local_nome: descricao.trim(),
      });
      Alert.alert("Sucesso", "Avistamento registrado! Ele já aparece no mapa.");
      navigation.goBack(); 
    } catch (err) {
      console.error("Erro ao salvar:", err);
      Alert.alert("Não foi possível registrar", mensagemErro(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView style={style.SafeAreaView}>
      
      {/* CABEÇALHO */}
      <View style={style.header}>
        <View style={style.headerSpacer} />
        
        <View style={style.headerTitleBox}>
          <Text style={style.headerTitleText}>NOVA OCORRÊNCIA</Text>
        </View>

        <TouchableOpacity onPress={() => navigation.goBack()} style={style.backButton}>
          <Ionicons name="arrow-back" size={24} color="#FFFFFF" />
        </TouchableOpacity>
      </View>

      <View style={style.formContainer}>
        <Text style={style.chipsLabel}>Animal avistado</Text>
        {erroAnimais ? (
          <TouchableOpacity onPress={carregarAnimais}>
            <Text style={style.chipsErro}>Não foi possível carregar os animais. Toque para tentar de novo.</Text>
          </TouchableOpacity>
        ) : (
          <ScrollView horizontal showsHorizontalScrollIndicator={false} style={style.chipsRow}>
            {animais.map((a) => (
              <TouchableOpacity
                key={a.id}
                style={[style.chip, animalId === a.id && style.chipAtivo]}
                onPress={() => setAnimalId(a.id)}
              >
                <Text style={style.chipText}>{a.nome_popular}</Text>
              </TouchableOpacity>
            ))}
          </ScrollView>
        )}

        <TextInput
          style={[style.input, style.textArea]}
          placeholder="Local ou descrição (até 120 caracteres)..."
          placeholderTextColor="#E0E0E0"
          value={descricao}
          onChangeText={setDescricao}
          maxLength={120}
          multiline
          numberOfLines={4}
        />

        <View style={style.locationCard}>
          <Text style={style.locationTitle}>Localização da Ocorrência</Text>
          {coordenadas ? (
            <Text style={style.locationCoords}>
              Lat: {coordenadas.latitude.toFixed(5)} {"\n"}Lon: {coordenadas.longitude.toFixed(5)}
            </Text>
          ) : (
            <Text style={style.locationCoords}>Nenhum local marcado</Text>
          )}
          
          <TouchableOpacity 
            style={style.mapButton} 
            onPress={abrirMapa}
            disabled={obtendoLocalizacao}
          >
            <Text style={style.mapButtonText}>
              {obtendoLocalizacao ? "Obtendo GPS..." : (coordenadas ? "Ajustar no Mapa" : "Marcar Localização")}
            </Text>
          </TouchableOpacity>
        </View>
      </View>


      <TouchableOpacity style={style.submitButton} onPress={salvarOcorrencia} disabled={loading}>
        <Text style={style.submitButtonText}>Gravar Ocorrência</Text>
      </TouchableOpacity>


      {loading && (
        <View style={style.centerContainer}>
          <ActivityIndicator size="large" color="#FFFFFF" />
          <Text style={style.loadingText}>Salvando registro...</Text>
        </View>
      )}


      <Modal visible={mapaVisivel} animationType="slide" transparent={false}>
        <View style={style.mapContainer}>
          {regiaoMapa && (
            <MapView
              style={style.map}
              initialRegion={regiaoMapa}
              onRegionChangeComplete={(region) => setRegiaoMapa(region)}
              showsUserLocation={true}
            />
          )}
          
          <View style={style.fixedMarker}>
            <View style={style.markerDot} />
            <View style={style.markerStick} />
          </View>

          <View style={style.mapFooter}>
            <TouchableOpacity style={style.cancelMapButton} onPress={() => setMapaVisivel(false)}>
              <Text style={style.cancelMapButtonText}>Cancelar</Text>
            </TouchableOpacity>
            
            <TouchableOpacity style={style.confirmMapButton} onPress={confirmarLocalizacao}>
              <Text style={style.confirmMapButtonText}>Confirmar Local</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </SafeAreaView>
  );
};